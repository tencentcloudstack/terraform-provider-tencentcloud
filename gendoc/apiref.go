package main

// apiref.go builds the "resource -> cloud API" mapping used to render the
// optional "API Reference" section of generated docs.
//
// The mapping is derived from two sources:
//
//  1. Static analysis (go/ast) of the provider source tree:
//     resource lifecycle function -> service method call -> (recursively)
//     further service methods / SDK client action calls
//     (e.g. `me.client.UseVpcClient().CreateVpc(request)`).
//     The SDK client method name is identical to the cloud API action name;
//     the client accessor (UseXxxClient) is resolved to a catalog product.
//
//  2. TencentCloud's international API catalog (fetched live by 'make
//     doc', with a bundled snapshot at gendoc/data/api-intl-snapshot.json
//     as the offline fallback), which resolves (product, action) -> doc
//     URL. Entries without a catalog URL fall back to the international
//     site's search URL for the action name.
//
// SDK import aliases do not always match the catalog product key, so an
// alias table built from the SDK import paths (plus a tiny manual override
// list) maps client accessor name -> catalog product.

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// apiRef describes a single cloud API used by a resource lifecycle phase.
// The rendered section lists only the action name and its doc URL — no
// descriptions.
type apiRef struct {
	Product string `json:"product"`
	Action  string `json:"action"`
	URL     string `json:"url"`
}

// manualProductAlias overrides the automatic product resolution for the few
// clients whose SDK product name differs from the international docs slug.
var manualProductAlias = map[string]string{
	"dasb":  "dasb", // placeholder; real aliases are resolved by action-set matching
	"css":   "css",
	"audit": "audit",
}

// serviceMethod describes one method declared on a service struct in the
// tencentcloud/services tree.
type serviceMethod struct {
	File string
	Recv string // receiver type name, e.g. "VpcService"
}

// apiAnalyzer performs the static analysis and catalog resolution.
type apiAnalyzer struct {
	fset         *token.FileSet
	repoRoot     string
	servicesRoot string

	// clientProd maps lowercase UseXxxClient accessor name -> catalog product.
	clientProd map[string]string
	// catalog resolves (product, action) -> international doc URL.
	catalog *catalog
	// clientActions records every (product -> set of actions) call found
	// in the services tree, used to align SDK product names with
	// international docs product slugs.
	clientActions map[string]map[string]bool
	// knownClients is the set of lowercase accessor names with a product.
	knownClients map[string]bool

	// methodIndex maps method name -> declared service methods (a name may
	// exist on several service structs, hence a list).
	methodIndex map[string][]serviceMethod
	// serviceCtors maps NewXxxService constructor name -> returned struct
	// type name. The key carries no package qualifier: this assumes that a
	// constructor name maps to exactly one type across all services
	// packages (true at the time of writing — 84 constructors, 0
	// collisions). If a same-name constructor ever returns different types
	// per package, type inference here must become package-aware.
	serviceCtors map[string]string

	// parsedFiles caches ASTs by path.
	parsedFiles map[string]*ast.File
}

func newAPIAnalyzer(repoRoot string, cat *catalog) (*apiAnalyzer, error) {
	clientProd, err := buildClientProductIndex(repoRoot)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		cat = loadCatalog()
	}
	known := map[string]bool{}
	for c := range clientProd {
		known[c] = true
	}
	a := &apiAnalyzer{
		fset:         token.NewFileSet(),
		repoRoot:     repoRoot,
		servicesRoot: filepath.Join(repoRoot, "tencentcloud", "services"),
		clientProd:   clientProd,
		catalog:      cat,
		knownClients: known,
		methodIndex:  map[string][]serviceMethod{},
		serviceCtors: map[string]string{},
		parsedFiles:  map[string]*ast.File{},
	}
	if err := a.indexServiceMethods(); err != nil {
		return nil, err
	}
	// Collect every (client product, action) call in the services tree so
	// the catalog can align SDK product names with international docs
	// products whose slugs differ (e.g. SDK "redis" vs docs slug "crs").
	a.collectClientActions()
	cat.alignCatalogWithSDKProducts(a.clientActions)
	return a, nil
}

// collectClientActions walks every indexed service method body and
// records every UseXxxClient().Action(...) invocation per product. The
// method index does not keep function names on its entries, so index
// them here first.
func (a *apiAnalyzer) collectClientActions() {
	a.clientActions = map[string]map[string]bool{}
	// Group entries by file so each parsed file is scanned once.
	byFile := map[string][]*ast.FuncDecl{}
	for name, matches := range a.methodIndex {
		for _, m := range matches {
			fd := a.lookupFuncDecl(m.File, name)
			if fd != nil {
				byFile[m.File] = append(byFile[m.File], fd)
			}
		}
	}
	for _, fds := range byFile {
		for _, fd := range fds {
			a.scanFuncForClientActions(fd)
		}
	}
}

// scanFuncForClientActions records client action calls in one function
// body (direct UseXxxClient().Action chains and locally held clients).
func (a *apiAnalyzer) scanFuncForClientActions(fd *ast.FuncDecl) {
	if fd.Body == nil {
		return
	}
	clientVars := map[string]string{}
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		inferClientVars(n, clientVars)
		return true
	})
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		useName := a.findClientUse(sel.X, clientVars)
		if useName == "" || !a.knownClients[useName] || !isAPIActionName(sel.Sel.Name) {
			return true
		}
		prod := a.clientProd[useName]
		if prod == "" {
			return true
		}
		if a.clientActions[prod] == nil {
			a.clientActions[prod] = map[string]bool{}
		}
		a.clientActions[prod][stripContextSuffix(sel.Sel.Name)] = true
		return true
	})
}

// indexServiceMethods walks the services tree once and records every method
// declared with a receiver, plus package-level helper functions that contain
// direct client calls (common in older resources). Test files are skipped.
func (a *apiAnalyzer) indexServiceMethods() error {
	return filepath.Walk(a.servicesRoot, func(path string, fi os.FileInfo, err error) error {
		if err != nil || fi.IsDir() || !strings.HasSuffix(path, ".go") ||
			strings.HasSuffix(path, "_test.go") {
			return nil
		}
		f, err := a.parseFile(path)
		if err != nil {
			return nil // unparseable files are simply not indexed
		}
		for _, d := range f.Decls {
			fd, ok := d.(*ast.FuncDecl)
			if !ok {
				continue
			}
			recv := ""
			if fd.Recv != nil && len(fd.Recv.List) > 0 {
				recv = recvTypeName(fd.Recv.List[0].Type)
			}
			if recv != "" {
				a.methodIndex[fd.Name.Name] = append(a.methodIndex[fd.Name.Name],
					serviceMethod{File: path, Recv: recv})
				continue
			}
			// Package-level function containing a direct client call:
			// index it too so handleCall can expand into it.
			hasClientCall := false
			ast.Inspect(fd, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				if sel, ok := call.Fun.(*ast.SelectorExpr); ok && isClientAccessor(sel.Sel.Name) {
					hasClientCall = true
					return false
				}
				return true
			})
			if hasClientCall {
				a.methodIndex[fd.Name.Name] = append(a.methodIndex[fd.Name.Name],
					serviceMethod{File: path, Recv: ""})
			}

			// NewXxxService(...) constructor: record the returned service
			// struct type so `x := NewXxxService(...)` type inference works.
			if strings.HasPrefix(fd.Name.Name, "New") && strings.HasSuffix(fd.Name.Name, "Service") &&
				fd.Type != nil && fd.Type.Results != nil && fd.Type.Results.NumFields() == 1 {
				rt := fd.Type.Results.List[0].Type
				if id, ok := rt.(*ast.Ident); ok {
					a.serviceCtors[fd.Name.Name] = id.Name
				}
			}
		}
		return nil
	})
}

func (a *apiAnalyzer) parseFile(path string) (*ast.File, error) {
	if f, ok := a.parsedFiles[path]; ok {
		return f, nil
	}
	f, err := parser.ParseFile(a.fset, path, nil, 0)
	if err != nil {
		return nil, err
	}
	a.parsedFiles[path] = f
	return f, nil
}

// analyzeLifecycle traces one lifecycle function body to the cloud APIs it
// invokes. It collects service method calls made on local variables whose
// concrete type is a service struct (e.g. `service := VpcService{...}`),
// then recursively expands service methods (following same-receiver calls
// and cross-service calls), collecting UseXxxClient().Action(...)
// invocations. The visited set guards against cycles.
func (a *apiAnalyzer) analyzeLifecycle(fn *ast.FuncDecl) []apiRef {
	seenRef := map[apiRef]bool{}
	visitedMethod := map[string]bool{} // "Recv.Method" -> expanded
	var refs []apiRef

	// localTypes infers variable -> service struct type from composite
	// literal assignments and var declarations in the function body.
	localTypes := map[string]string{}
	if fn.Body != nil {
		ast.Inspect(fn.Body, func(n ast.Node) bool {
			inferLocalTypes(n, localTypes, a.serviceCtors)
			return true
		})
	}

	var walkFuncDecl func(fd *ast.FuncDecl, localTypes map[string]string)

	// handleCall inspects one call expression. Returns the recv hint used
	// for nested calls when the call expands into a service method, or "".
	handleCall := func(call *ast.CallExpr, hint string, localTypes map[string]string, clientVars map[string]string) string {
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			// Plain function call (e.g. package-level helper like
			// haVipAssociateEip or a goroutine delegate doResource...Create).
			if id, ok := call.Fun.(*ast.Ident); ok && !visitedMethod[""+id.Name] {
				if matches, ok := a.methodIndex[id.Name]; ok && len(matches) > 0 {
					pick := &matches[0]
					visitedMethod[""+id.Name] = true
					if fd := a.lookupFuncDecl(pick.File, id.Name); fd != nil {
						walkFuncDecl(fd, nil)
					}
				}
			}
			return ""
		}
		action := sel.Sel.Name

		// SDK client action call? (receiver chain contains UseXxxClient())
		if useName := a.findClientUse(sel.X, clientVars); useName != "" && a.knownClients[useName] {
			if !isAPIActionName(action) {
				return ""
			}
			if ref := a.resolveAPI(apiCallKey{client: useName, action: stripContextSuffix(action)}); ref != nil && !seenRef[*ref] {
				seenRef[*ref] = true
				refs = append(refs, *ref)
			}
			return ""
		}

		// Determine the expected receiver type: from the hint (receiver of
		// the enclosing method), or from a local service variable.
		recvHint := hint
		if id, ok := sel.X.(*ast.Ident); ok {
			if t, ok := localTypes[id.Name]; ok {
				recvHint = t
			}
		}

		// Service method call? Only methods present in the method index.
		matches, ok := a.methodIndex[action]
		if !ok {
			return ""
		}
		// Disambiguate by receiver type. When the call site's receiver is
		// unknown (recvHint empty) and the name is declared on several
		// distinct service structs, the target is ambiguous: skip it rather
		// than guessing and emitting a wrong-product API link. A single
		// declaration is always safe.
		var pick *serviceMethod
		for i := range matches {
			if matches[i].Recv == recvHint {
				pick = &matches[i]
				break
			}
		}
		if pick == nil {
			if recvHint == "" && distinctRecvs(matches) > 1 {
				return ""
			}
			pick = &matches[0]
		}
		key := pick.Recv + "." + action
		if visitedMethod[key] {
			return ""
		}
		visitedMethod[key] = true
		if fd := a.lookupFuncDecl(pick.File, action); fd != nil {
			walkFuncDecl(fd, nil)
		}
		return pick.Recv
	}

	walkFuncDecl = func(fd *ast.FuncDecl, types map[string]string) {
		recv := ""
		if fd.Recv != nil && len(fd.Recv.List) > 0 {
			recv = recvTypeName(fd.Recv.List[0].Type)
		}
		lt := types
		if lt == nil {
			lt = map[string]string{}
		}
		clientVars := map[string]string{}
		if fd.Body != nil {
			// Two passes: first collect local variable types (service
			// structs and held clients), then walk the call tree with the
			// completed maps.
			ast.Inspect(fd.Body, func(n ast.Node) bool {
				inferLocalTypes(n, lt, a.serviceCtors)
				inferClientVars(n, clientVars)
				return true
			})
			ast.Inspect(fd.Body, func(n ast.Node) bool {
				if call, ok := n.(*ast.CallExpr); ok {
					handleCall(call, recv, lt, clientVars)
				}
				return true // never prune; recursion is guarded by visited sets
			})
		}
	}

	walkFuncDecl(fn, localTypes)
	sort.Slice(refs, func(i, j int) bool { return refs[i].Action < refs[j].Action })
	return refs
}

// lookupFuncDecl finds the function/method declaration named fnName in file.
func (a *apiAnalyzer) lookupFuncDecl(file, fnName string) *ast.FuncDecl {
	f, err := a.parseFile(file)
	if err != nil {
		return nil
	}
	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok && fd.Name.Name == fnName {
			return fd
		}
	}
	return nil
}

// lookupMethodByReceiver finds the method fnName declared on receiver type
// recv (searching files that declare methods on that receiver). Falls back
// to the first declaration found.
func (a *apiAnalyzer) lookupMethodByReceiver(recv, fnName string) *ast.FuncDecl {
	matches, ok := a.methodIndex[fnName]
	if !ok {
		return nil
	}
	for _, m := range matches {
		if m.Recv == recv {
			if fd := a.lookupFuncDecl(m.File, fnName); fd != nil {
				return fd
			}
		}
	}
	if len(matches) > 0 {
		return a.lookupFuncDecl(matches[0].File, fnName)
	}
	return nil
}

// findClientUse walks an expression tree looking for a UseXxxClient() call
// and returns the lowercase client accessor name, or "" if none found.
// clientVars maps local variable names holding a client (e.g. `client :=
// me.client.UseGaapClient()`) to their accessor name.
func (a *apiAnalyzer) findClientUse(expr ast.Expr, clientVars map[string]string) string {
	switch e := expr.(type) {
	case *ast.CallExpr:
		if sel, ok := e.Fun.(*ast.SelectorExpr); ok && isClientAccessor(sel.Sel.Name) {
			name := strings.TrimSuffix(strings.TrimPrefix(sel.Sel.Name, "Use"), "Client")
			if name != "" {
				return strings.ToLower(name)
			}
		}
		return a.findClientUse(e.Fun, clientVars)
	case *ast.SelectorExpr:
		return a.findClientUse(e.X, clientVars)
	case *ast.ParenExpr:
		return a.findClientUse(e.X, clientVars)
	case *ast.Ident:
		// A local variable holding a client, used later as
		// `client.Action(...)`.
		return clientVars[e.Name]
	}
	return ""
}

func isClientAccessor(name string) bool {
	return strings.HasPrefix(name, "Use") && strings.HasSuffix(name, "Client") && len(name) > len("UseClient")
}

// isAPIActionName filters out wrapper methods and unexported helpers so
// only real SDK API invocations remain. SDK methods suffixed with
// "WithContext" are accepted with the suffix stripped (the underlying API
// action is identical).
func isAPIActionName(name string) bool {
	if !ast.IsExported(name) {
		return false
	}
	switch {
	case strings.HasSuffix(name, "WithRetry"), strings.HasSuffix(name, "AwaitTask"),
		name == "WithHttpTransport":
		return false
	case strings.HasSuffix(name, "WithContext"):
		return len(strings.TrimSuffix(name, "WithContext")) > 0
	}
	return true
}

// stripContextSuffix removes the SDK "WithContext" wrapper suffix.
func stripContextSuffix(name string) string {
	return strings.TrimSuffix(name, "WithContext")
}

// apiCallKey identifies one client action invocation site.
type apiCallKey struct {
	client string // lowercase client accessor name (e.g. "vpc")
	action string // SDK method == cloud API action (e.g. "CreateVpc")
}

// intlSearchURL is the fallback for actions without a catalogue doc URL:
// the international documentation site's search page for the action name.
const intlSearchURL = "https://www.tencentcloud.com/search/?wd=%s&type=document&lang=en"

// resolveAPI turns a call key into an international-catalogue-backed
// apiRef. The doc URL comes from the international docs catalogue keyed
// by (product, action); entries without a page fall back to the
// international site's search page. Every rendered link is
// international — no Chinese-site URLs are involved. Returns nil when
// the (client, action) pair has no catalogue page and the action name
// does not look like a real API action (the catalogue doubles as a
// whitelist against false positives when available).
func (a *apiAnalyzer) resolveAPI(k apiCallKey) *apiRef {
	prod := a.clientProd[k.client]
	if prod == "" {
		return nil
	}
	if url := a.catalog.indexForClient(prod, k.action); url != "" {
		return &apiRef{Product: prod, Action: k.action, URL: url}
	}
	// No international page: degrade to the site's search URL for the
	// action name.
	return &apiRef{
		Product: prod,
		Action:  k.action,
		URL:     fmt.Sprintf(intlSearchURL, k.action),
	}
}

// distinctRecvs counts the distinct receiver type names among matches.
func distinctRecvs(matches []serviceMethod) int {
	set := map[string]bool{}
	for _, m := range matches {
		set[m.Recv] = true
	}
	return len(set)
}

// recvTypeName extracts the receiver type name from a receiver field type
// (handles *T and T).
func recvTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return recvTypeName(t.X)
	case *ast.Ident:
		return t.Name
	case *ast.IndexExpr:
		return recvTypeName(t.X)
	}
	return ""
}

// buildClientProductIndex maps every UseXxxClient accessor to its catalog
// product using the SDK import path of the client's declared return type.
func buildClientProductIndex(repoRoot string) (map[string]string, error) {
	clientFile := filepath.Join(repoRoot, "tencentcloud", "connectivity", "client.go")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, clientFile, nil, parser.ImportsOnly)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", clientFile, err)
	}

	// alias -> SDK product directory from the import path.
	importDirs := map[string]string{}
	for _, imp := range f.Imports {
		path := strings.Trim(imp.Path.Value, `"`)
		// Also match the vendored intl SDK (tencentcloud-sdk-go-intl-en):
		// a few resources call the international-site clients through it.
		// Both SDKs are vendored in-repo; no external download is involved.
		if !strings.Contains(path, "tencentcloud-sdk-go/tencentcloud/") &&
			!strings.Contains(path, "tencentcloud-sdk-go-intl-en/tencentcloud/") {
			continue
		}
		alias := filepath.Base(path)
		if imp.Name != nil {
			alias = imp.Name.Name
		}
		// .../tencentcloud/<product>/v<version>
		parts := strings.Split(path, "/")
		if len(parts) >= 2 {
			importDirs[strings.ToLower(alias)] = parts[len(parts)-2]
		}
	}

	// Scan UseXxxClient declarations to learn the concrete return type
	// (package alias) of each client accessor.
	pkgs, err := parser.ParseDir(fset, filepath.Dir(clientFile), nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parse dir %s: %w", filepath.Dir(clientFile), err)
	}
	clientPkg := map[string]string{} // lowercase accessor name -> import alias
	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				fd, ok := decl.(*ast.FuncDecl)
				if !ok || fd.Recv == nil || !isClientAccessor(fd.Name.Name) {
					continue
				}
				name := strings.TrimSuffix(strings.TrimPrefix(fd.Name.Name, "Use"), "Client")
				res := fd.Type.Results
				if res == nil || len(res.List) != 1 {
					continue
				}
				star, ok := res.List[0].Type.(*ast.StarExpr)
				if !ok {
					continue
				}
				sel, ok := star.X.(*ast.SelectorExpr)
				if !ok {
					continue
				}
				ident, ok := sel.X.(*ast.Ident)
				if !ok {
					continue
				}
				clientPkg[strings.ToLower(name)] = ident.Name
			}
		}
	}

	// Resolve accessor name -> catalog product, stripping trailing version
	// digits (cvmv20170312 -> cvm).
	idx := map[string]string{}
	for use, alias := range clientPkg {
		dir := importDirs[strings.ToLower(alias)]
		if dir == "" {
			continue
		}
		idx[use] = strings.ToLower(dir)
	}
	// Manual overrides last, so they win.
	for k, v := range manualProductAlias {
		idx[k] = v
	}
	return idx, nil
}

// renderAPIRefsFor renders the markdown bullet list for a phase -> refs
// map, with the lead-in sentence matching the doc type.
func renderAPIRefsFor(dtype string, refs map[string][]apiRef) string {
	leads := map[string]string{
		"resource":           "This resource interacts with the following TencentCloud APIs:",
		"data_source":        "This data source queries the following TencentCloud APIs:",
		"action":             "This action invokes the following TencentCloud APIs:",
		"ephemeral":          "This ephemeral resource reads the following TencentCloud APIs:",
		"ephemeral_resource": "This ephemeral resource reads the following TencentCloud APIs:",
	}
	lead, ok := leads[dtype]
	if !ok {
		lead = leads["resource"]
	}
	return renderAPIRefsLead(lead, dtype, refs)
}

// renderAPIRefs renders the markdown bullet list for a phase -> refs map
// with the resource lead-in.
func renderAPIRefs(refs map[string][]apiRef) string {
	return renderAPIRefsLead("This resource interacts with the following TencentCloud APIs:", "resource", refs)
}

// renderAPIRefsLead is the shared renderer. verbOrder controls the phase
// sequence; verbs outside the map sort after the known ones alphabetically.
func renderAPIRefsLead(lead string, dtype string, refs map[string][]apiRef) string {
	nonEmpty := 0
	for _, list := range refs {
		nonEmpty += len(list)
	}
	if nonEmpty == 0 {
		return ""
	}
	verbOrder := map[string]int{"Create": 0, "Read": 1, "Update": 2, "Delete": 3, "Invoke": 4, "Open": 5}
	var b strings.Builder
	b.WriteString(lead + "\n\n")
	verbs := make([]string, 0, len(refs))
	for v := range refs {
		verbs = append(verbs, v)
	}
	sort.Slice(verbs, func(i, j int) bool { return verbOrder[verbs[i]] < verbOrder[verbs[j]] })
	for _, v := range verbs {
		list := refs[v]
		if len(list) == 0 {
			continue
		}
		sort.Slice(list, func(i, j int) bool { return list[i].Action < list[j].Action })
		for _, r := range list {
			if r.URL != "" {
				fmt.Fprintf(&b, "* %s: [%s](%s)\n", v, r.Action, r.URL)
			} else {
				fmt.Fprintf(&b, "* %s: `%s`\n", v, r.Action)
			}
		}
	}
	return strings.TrimRight(b.String(), "\n")
}

// ctorReturnType resolves the service struct type returned by a
// NewXxxService(...) call, supporting both plain and package-qualified
// references (NewRedisService(...) and svccrs.NewRedisService(...)).
func ctorReturnType(call *ast.CallExpr, ctors map[string]string) (string, string) {
	switch fn := call.Fun.(type) {
	case *ast.Ident:
		if t, ok := ctors[fn.Name]; ok {
			return t, fn.Name
		}
	case *ast.SelectorExpr:
		if t, ok := ctors[fn.Sel.Name]; ok {
			return t, fn.Sel.Name
		}
	}
	return "", ""
}

// inferLocalTypes records variable -> struct-type bindings from composite
// literal assignments (x := Service{...}, var x = Service{...}) and
// NewXxxService(...) constructor calls so that service method receivers can
// be resolved precisely. ctors maps constructor name -> service type and may
// be nil.
func inferLocalTypes(n ast.Node, out map[string]string, ctors map[string]string) {
	switch e := n.(type) {
	case *ast.AssignStmt:
		if len(e.Lhs) != len(e.Rhs) {
			return
		}
		for i, rhs := range e.Rhs {
			if call, ok := rhs.(*ast.CallExpr); ok && ctors != nil {
				if t, fnName := ctorReturnType(call, ctors); fnName != "" {
					if lhs, ok := e.Lhs[i].(*ast.Ident); ok {
						out[lhs.Name] = t
					}
				}
			}
			lit, ok := rhs.(*ast.CompositeLit)
			if !ok {
				continue
			}
			if id, ok := lit.Type.(*ast.Ident); ok {
				if lhs, ok := e.Lhs[i].(*ast.Ident); ok {
					out[lhs.Name] = id.Name
				}
			}
		}
	case *ast.DeclStmt:
		gd, ok := e.Decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			return
		}
		for _, spec := range gd.Specs {
			ps, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}
			// var x ServiceType{...} — type on the spec.
			if id, ok := ps.Type.(*ast.Ident); ok {
				for _, name := range ps.Names {
					out[name.Name] = id.Name
				}
				continue
			}
			// var x = ServiceType{...} or x = NewXxxService(...) — type
			// from the composite literal or constructor.
			if len(ps.Names) == len(ps.Values) {
				for i, val := range ps.Values {
					if call, ok := val.(*ast.CallExpr); ok && ctors != nil {
						if t, fnName := ctorReturnType(call, ctors); fnName != "" {
							out[ps.Names[i].Name] = t
						}
					}
					lit, ok := val.(*ast.CompositeLit)
					if !ok {
						continue
					}
					if id, ok := lit.Type.(*ast.Ident); ok {
						out[ps.Names[i].Name] = id.Name
					}
				}
			}
		}
	}
}

// inferClientVars records local variables initialised directly from a
// UseXxxClient() call (x := me.client.UseGaapClient()) so later plain
// `x.Action(...)` calls resolve to the client.
func inferClientVars(n ast.Node, out map[string]string) {
	record := func(lhs *ast.Ident, rhs ast.Expr) {
		if call, ok := rhs.(*ast.CallExpr); ok {
			if sel, ok := call.Fun.(*ast.SelectorExpr); ok && isClientAccessor(sel.Sel.Name) {
				name := strings.TrimSuffix(strings.TrimPrefix(sel.Sel.Name, "Use"), "Client")
				if name != "" {
					out[lhs.Name] = strings.ToLower(name)
				}
			}
		}
	}
	switch e := n.(type) {
	case *ast.AssignStmt:
		if len(e.Lhs) == len(e.Rhs) {
			for i, rhs := range e.Rhs {
				if lhs, ok := e.Lhs[i].(*ast.Ident); ok {
					record(lhs, rhs)
				}
			}
		}
	case *ast.DeclStmt:
		gd, ok := e.Decl.(*ast.GenDecl)
		if !ok || gd.Tok != token.VAR {
			return
		}
		for _, spec := range gd.Specs {
			ps, ok := spec.(*ast.ValueSpec)
			if !ok || len(ps.Names) != len(ps.Values) {
				continue
			}
			for i, val := range ps.Values {
				record(ps.Names[i], val)
			}
		}
	}
}

// ---- official generation pipeline integration ----

// sharedAPIAnalyzer lazily initialised analyzer for the doc pipelines.
var sharedAPIAnalyzer *apiAnalyzer

// getSharedAPIAnalyzer returns the process-wide analyzer, initialising it on
// first use. Returns nil when analysis is unavailable (e.g. running outside
// the repo root); callers must handle nil.
func getSharedAPIAnalyzer() *apiAnalyzer {
	if sharedAPIAnalyzer != nil {
		return sharedAPIAnalyzer
	}
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		return nil
	}
	cat := loadCatalog()
	a, err := newAPIAnalyzer(repoRoot, cat)
	if err != nil {
		message("[WARN]API reference analysis unavailable: %v\n", err)
		return nil
	}
	sharedAPIAnalyzer = a
	return a
}

// renderResourceAPISection produces the "API Reference" markdown body for
// one SDKv2 resource. Order of precedence:
//
//  1. a manual "API" section in the resource .md (authoritative override);
//  2. static analysis + catalog lookup.
func renderResourceAPISection(analyzer *apiAnalyzer, repoRoot, resName string, res *schema.Resource, manual string) string {
	manual = strings.TrimSpace(manual)
	if manual != "" {
		return manual
	}
	if analyzer == nil || res == nil {
		return ""
	}
	resFile := findResourceFile(repoRoot, resName)
	if resFile == "" {
		return ""
	}
	refsByVerb := map[string][]apiRef{}
	for verb, fnPtr := range map[string]func(*schema.ResourceData, interface{}) error{
		"Create": res.Create, "Read": res.Read, "Update": res.Update, "Delete": res.Delete,
	} {
		if fnPtr == nil {
			continue
		}
		fd := findFuncDeclByPointer(analyzer, resFile, fnPtr)
		if fd == nil {
			continue
		}
		refsByVerb[verb] = analyzer.analyzeLifecycle(fd)
	}
	return renderAPIRefs(dedupVerbRefs(refsByVerb))
}

// dedupVerbRefs dedups and sorts refs per verb.
func dedupVerbRefs(refs map[string][]apiRef) map[string][]apiRef {
	out := map[string][]apiRef{}
	for verb, list := range refs {
		seen := map[apiRef]bool{}
		var uniq []apiRef
		for _, r := range list {
			if !seen[r] {
				seen[r] = true
				uniq = append(uniq, r)
			}
		}
		out[verb] = uniq
	}
	return out
}

// findFuncDeclByPointer resolves the runtime function pointer of a lifecycle
// function back to its declaring FuncDecl via the debug symbol table.
func findFuncDeclByPointer(a *apiAnalyzer, resFile string, fn func(*schema.ResourceData, interface{}) error) *ast.FuncDecl {
	name := runtimeFuncName(fn)
	if name == "" {
		return nil
	}
	// The function may live in the resource file or a sibling file.
	if fd := a.lookupFuncDecl(resFile, name); fd != nil {
		return fd
	}
	// Search same directory.
	dir := filepath.Dir(resFile)
	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") || strings.HasSuffix(e.Name(), "_test.go") {
			continue
		}
		if fd := a.lookupFuncDecl(filepath.Join(dir, e.Name()), name); fd != nil {
			return fd
		}
	}
	return nil
}

// frameworkResourceRecvType resolves the Go struct receiver type name of the
// framework resource implementation (e.g. "SsmSecretVersionV2Resource" for
// tencentcloud_ssm_secret_version_v2) by scanning the declaring file for a
// method named "Metadata" and reading its receiver.
func frameworkResourceRecvType(analyzer *apiAnalyzer, resFile, resName string) string {
	f, err := analyzer.parseFile(resFile)
	if err != nil {
		return ""
	}
	for _, d := range f.Decls {
		fd, ok := d.(*ast.FuncDecl)
		if !ok || fd.Name.Name != "Metadata" || fd.Recv == nil || len(fd.Recv.List) == 0 {
			continue
		}
		if recv := recvTypeName(fd.Recv.List[0].Type); recv != "" {
			return recv
		}
	}
	// Fallback: derive from the type name embedded in the file name.
	base := filepath.Base(resFile)
	parts := strings.Split(strings.TrimSuffix(base, ".go"), "_")
	for i := len(parts) - 1; i > 0; i-- {
		if parts[i-1] == "tc" {
			cand := strings.Join(parts[i:], "")
			cand = strings.ToUpper(cand[:1]) + cand[1:] + "Resource"
			return cand
		}
	}
	return ""
}

// ---- datasource / action / ephemeral rendering ----

// frameworkRefKind describes a framework reference type whose lifecycle
// methods the API section should trace.
type frameworkRefKind struct {
	// verbs maps the rendered verb to the Go method name on the impl.
	verbs map[string]string
}

// frameworkRefKinds lists the lifecycle verbs per framework reference type.
// Actions implement Invoke; ephemeral resources implement Open.
var frameworkRefKinds = map[string]frameworkRefKind{
	"action":             {verbs: map[string]string{"Invoke": "Invoke"}},
	"ephemeral_resource": {verbs: map[string]string{"Open": "Open"}},
}

// renderFrameworkRefAPISection produces the "API Reference" body for
// framework actions and ephemeral resources. The receiver type is resolved
// from the declaring file via its Metadata method.
func renderFrameworkRefAPISection(analyzer *apiAnalyzer, repoRoot, name, kind string) string {
	if analyzer == nil {
		return ""
	}
	spec, ok := frameworkRefKinds[kind]
	if !ok {
		return ""
	}
	resFile := findResourceFile(repoRoot, name)
	if resFile == "" {
		return ""
	}
	recv := frameworkResourceRecvType(analyzer, resFile, name)
	if recv == "" {
		return ""
	}
	refsByVerb := map[string][]apiRef{}
	for verb, method := range spec.verbs {
		fd := analyzer.lookupMethodByReceiver(recv, method)
		if fd != nil {
			refsByVerb[verb] = analyzer.analyzeLifecycle(fd)
		}
	}
	return renderAPIRefsFor(kind, dedupVerbRefs(refsByVerb))
}

// readDataSourceAPIs traces one SDKv2 data source's Read function.
func (a *apiAnalyzer) readDataSourceAPIs(repoRoot, name string, ds *schema.Resource) []apiRef {
	if a == nil || ds == nil || ds.Read == nil {
		return nil
	}
	resFile := findResourceFile(repoRoot, name)
	if resFile == "" {
		return nil
	}
	fd := findFuncDeclByPointer(a, resFile, ds.Read)
	if fd == nil {
		return nil
	}
	return a.analyzeLifecycle(fd)
}

// balanceFences guards against source .md files with an odd number of
// code-fence markers. Rather than deleting content, the LAST lone fence
// marker is escaped (\`\`\`) so it renders as literal text: a marker with
// a language tag (```hcl) is a legitimate opener and is never touched,
// while an untagged lone marker is the stray one upstream left behind.
func balanceFences(s string) string {
	if strings.Count(s, "```")%2 == 0 {
		return s
	}
	lines := strings.Split(s, "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.TrimSpace(lines[i]) == "```" {
			lines[i] = "\\`" + "\\`" + "\\`"
			return strings.Join(lines, "\n")
		}
	}
	return s
}

// inlineFenceRe matches code-fence markers embedded inside a schema
// description (upstream pasted table markup like ":```open: ...").
var inlineFenceRe = regexp.MustCompile("`{3,}")

// sanitizeInlineFences escapes stray code-fence markers inside a schema
// description so they cannot break the fenced blocks of the generated doc.
func sanitizeInlineFences(s string) string {
	if !strings.Contains(s, "```") {
		return s
	}
	return inlineFenceRe.ReplaceAllString(s, "\\`\\`\\`")
}
