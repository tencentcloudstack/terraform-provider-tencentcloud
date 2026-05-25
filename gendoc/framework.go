// Package main — framework.go is the framework-stack counterpart of the
// sdkv2 documentation generator. It enumerates the six framework
// reference types (Resource / DataSource / Function / EphemeralResource /
// ListResource / Action), reads their per-type markdown files placed
// next to the Go source under tencentcloud/framework/<product>/<type>/,
// and renders them into website/docs/<dir>/<resource>.html.markdown.
//
// Output dirs:
//
//	resource           -> website/docs/r/
//	datasource         -> website/docs/d/
//	function           -> website/docs/functions/
//	ephemeral resource -> website/docs/ephemeral-resources/
//	list resource      -> website/docs/list-resources/
//	action             -> website/docs/actions/
//
// The framework index file is tencentcloud/framework/provider.md. It
// follows the same "Resources List" syntax as the sdkv2 provider.md but
// adds four extra section headers — "Function", "Ephemeral Resource",
// "List Resource", "Action" — handled by GetFrameworkIndex.
package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	frameworkaction "github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	frameworklist "github.com/hashicorp/terraform-plugin-framework/list"
	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	cloud "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud"
	tcfw "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/framework"
)

// fwDocType is the discriminator for the six framework reference types.
type fwDocType string

const (
	fwResource  fwDocType = "resource"
	fwDataSrc   fwDocType = "data_source"
	fwFunction  fwDocType = "function"
	fwEphemeral fwDocType = "ephemeral_resource"
	fwList      fwDocType = "list_resource"
	fwAction    fwDocType = "action"
)

// outputDir maps a framework doc type to its registry website directory.
func (t fwDocType) outputDir() string {
	switch t {
	case fwResource:
		return "r"
	case fwDataSrc:
		return "d"
	case fwFunction:
		return "functions"
	case fwEphemeral:
		return "ephemeral-resources"
	case fwList:
		return "list-resources"
	case fwAction:
		return "actions"
	}
	return ""
}

// sourceMdRelDir maps a framework doc type to the in-tree path under
// tencentcloud/framework/<product>/ where the .md file is colocated with
// the Go source. The full relative path looks like:
//
//	framework/<product>/<sourceMdRelDir>/<resource>.md
//
// The <product> segment comes from the framework registry — meta for
// cross-product references, cvm for CVM, etc.
func (t fwDocType) sourceMdRelDir() string {
	switch t {
	case fwResource:
		return "resources"
	case fwDataSrc:
		return "datasources"
	case fwFunction:
		return "functions"
	case fwEphemeral:
		return "ephemerals"
	case fwList:
		return "lists"
	case fwAction:
		return "actions"
	}
	return ""
}

// frameworkProduct mirrors the sdkv2 Product struct but with the four
// framework-only collections appended.
type frameworkProduct struct {
	Name        string
	DataSources []string
	Resources   []string
	Functions   []string
	Ephemerals  []string
	Lists       []string
	Actions     []string
}

// genFrameworkDocs is the framework-side equivalent of the sdkv2 main
// loop. It is invoked from main() after the sdkv2 generation completes.
func genFrameworkDocs(repoRoot string) []frameworkProduct {
	frameworkRoot := filepath.Join(repoRoot, "framework")

	// 1) Build framework provider so we can drive 6 factory aggregators.
	primary := cloud.Provider()
	prov := tcfw.NewProvider(primary)
	ctx := context.Background()

	// 2) Read framework provider.md to learn the per-product groupings.
	products := readFrameworkIndex(frameworkRoot)

	// 3) Render each reference type.
	for _, fac := range prov.Resources(ctx) {
		r := fac()
		name := frameworkResourceTypeName(ctx, primary, r)
		genFrameworkDoc(frameworkRoot, fwResource, name, productOf(products, name, fwResource), func(out *map[string]fwAttrSpec) string {
			resp := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &resp)
			*out = flattenResourceSchema(resp.Schema)
			return strings.TrimSpace(resp.Schema.Description)
		})
	}
	for _, fac := range prov.DataSources(ctx) {
		ds := fac()
		name := frameworkDataSourceTypeName(ctx, primary, ds)
		genFrameworkDoc(frameworkRoot, fwDataSrc, name, productOf(products, name, fwDataSrc), func(out *map[string]fwAttrSpec) string {
			resp := datasource.SchemaResponse{}
			ds.Schema(ctx, datasource.SchemaRequest{}, &resp)
			*out = flattenDataSourceSchema(resp.Schema)
			return strings.TrimSpace(resp.Schema.Description)
		})
	}
	if pf, ok := prov.(frameworkprovider.ProviderWithFunctions); ok {
		for _, fac := range pf.Functions(ctx) {
			fn := fac()
			name := frameworkFunctionName(ctx, primary, fn)
			genFrameworkDoc(frameworkRoot, fwFunction, name, productOf(products, name, fwFunction), func(out *map[string]fwAttrSpec) string {
				resp := function.DefinitionResponse{}
				fn.Definition(ctx, function.DefinitionRequest{}, &resp)
				*out = flattenFunctionDefinition(resp.Definition)
				return strings.TrimSpace(firstNonEmpty(resp.Definition.MarkdownDescription, resp.Definition.Description, resp.Definition.Summary))
			})
		}
	}
	if pe, ok := prov.(frameworkprovider.ProviderWithEphemeralResources); ok {
		for _, fac := range pe.EphemeralResources(ctx) {
			er := fac()
			name := frameworkEphemeralName(ctx, primary, er)
			genFrameworkDoc(frameworkRoot, fwEphemeral, name, productOf(products, name, fwEphemeral), func(out *map[string]fwAttrSpec) string {
				resp := ephemeral.SchemaResponse{}
				er.Schema(ctx, ephemeral.SchemaRequest{}, &resp)
				*out = flattenEphemeralSchema(resp.Schema)
				return strings.TrimSpace(resp.Schema.Description)
			})
		}
	}
	// list resources are deliberately not enumerated from the registry
	// (the framework v1.19 list.ListResource interface needs a companion
	// managed resource that does not yet exist). The .md files placed
	// under framework/meta/lists/ are however parsed via the framework
	// provider.md index so that, once the real list reference is wired
	// up, the documentation pipeline keeps working with no changes.
	for _, name := range allListNames(products) {
		genFrameworkListPlaceholder(frameworkRoot, name, productOfList(products, name))
	}
	_ = frameworklist.ListResource(nil) // keep the import alive for future use
	if pa, ok := prov.(frameworkprovider.ProviderWithActions); ok {
		for _, fac := range pa.Actions(ctx) {
			a := fac()
			name := frameworkActionName(ctx, primary, a)
			genFrameworkDoc(frameworkRoot, fwAction, name, productOf(products, name, fwAction), func(out *map[string]fwAttrSpec) string {
				resp := frameworkaction.SchemaResponse{}
				a.Schema(ctx, frameworkaction.SchemaRequest{}, &resp)
				*out = flattenActionSchema(resp.Schema)
				return strings.TrimSpace(resp.Schema.Description)
			})
		}
	}
	return products
}

// schemaExtractor is the closure passed to genFrameworkDoc to delay the
// per-type-specific schema extraction until product/path lookups are
// done. It populates *out and returns the schema-level description.
type schemaExtractor func(out *map[string]fwAttrSpec) string

// genFrameworkDoc renders a single .md file for any framework reference
// type. It mirrors the sdkv2 genDoc function but reads the inline .md
// from framework/<product>/<type>/<name>.md and writes the rendered
// markdown to website/docs/<outputDir>/<name>.html.markdown.
func genFrameworkDoc(frameworkRoot string, dtype fwDocType, name, product string, extract schemaExtractor) {
	if name == "" {
		fail(fmt.Sprintf("framework %s has empty type name", dtype))
	}
	if !strings.HasPrefix(name, cloudPrefix) {
		fail(fmt.Sprintf("framework %s %q must start with %q", dtype, name, cloudPrefix))
	}
	resName := name[len(cloudPrefix):]

	// Locate the colocated .md file. The product string we read from the
	// framework provider.md is purely a UI label — the on-disk product
	// segment is fixed by the registry. We try both: explicit product
	// (lower-cased without parenthesised aliases) first, then a generic
	// "meta" fallback.
	mdPath := lookupFrameworkMd(frameworkRoot, dtype, resName, product)
	if mdPath == "" {
		fail(fmt.Sprintf("framework %s %q is missing its .md file under framework/<product>/%s/", dtype, name, dtype.sourceMdRelDir()))
	}

	raw, err := os.ReadFile(mdPath)
	if err != nil {
		fail(fmt.Sprintf("read %s failed: %s", mdPath, err))
	}
	desc := strings.TrimSpace(string(raw))
	if desc == "" {
		fail(fmt.Sprintf("description empty: %s", mdPath))
	}

	// Optional Import section (resource only).
	importBlock := ""
	if i := strings.Index(desc, "\nImport\n"); i != -1 {
		importBlock = strings.TrimSpace(desc[i+8:])
		desc = strings.TrimSpace(desc[:i])
	}

	// Required Example Usage section.
	example := ""
	if i := strings.Index(desc, "\nExample Usage\n"); i != -1 {
		example = formatHCL(desc[i+15:])
		desc = strings.TrimSpace(desc[:i])
	} else {
		fail(fmt.Sprintf("example usage missing: %s", mdPath))
	}

	descShort := desc
	if i := strings.Index(desc, "\n\n"); i != -1 {
		descShort = strings.TrimSpace(desc[:i])
	}

	attrs := map[string]fwAttrSpec{}
	_ = extract(&attrs)

	arguments, attributes := renderFrameworkSchema(name, attrs)
	if dtype == fwResource && !strings.Contains(attributes, "* `id`") {
		attributes = "* `id` - ID of the resource.\n" + attributes
	}

	out := map[string]string{
		"product":           product,
		"name":              name,
		"dtype":             string(dtype),
		"resource":          resName,
		"cloud_mark":        cloudMark,
		"cloud_title":       cloudTitle,
		"example":           example,
		"description":       desc,
		"description_short": descShort,
		"arguments":         arguments,
		"attributes":        attributes,
		"timeouts":          "",
		"import":            importBlock,
	}

	outPath := filepath.Join(docRoot, dtype.outputDir(), fmt.Sprintf("%s.html.markdown", resName))
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		fail(fmt.Sprintf("mkdir %s failed: %s", filepath.Dir(outPath), err))
	}
	fd, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		fail(fmt.Sprintf("open %s failed: %s", outPath, err))
	}
	defer fd.Close()

	t := template.Must(template.New("fw").Parse(docTPL))
	if err := t.Execute(fd, out); err != nil {
		fail(fmt.Sprintf("write %s failed: %s", outPath, err))
	}
	message("[SUCC.]write doc to file success: %s", outPath)
}

// genFrameworkListPlaceholder renders a list resource doc strictly from
// its colocated .md (no Go schema source available yet).
func genFrameworkListPlaceholder(frameworkRoot, name, product string) {
	resName := strings.TrimPrefix(name, cloudPrefix)
	mdPath := lookupFrameworkMd(frameworkRoot, fwList, resName, product)
	if mdPath == "" {
		fail(fmt.Sprintf("framework list %q is missing its .md file under framework/<product>/lists/", name))
	}
	raw, err := os.ReadFile(mdPath)
	if err != nil {
		fail(fmt.Sprintf("read %s failed: %s", mdPath, err))
	}
	desc := strings.TrimSpace(string(raw))
	if desc == "" {
		fail(fmt.Sprintf("description empty: %s", mdPath))
	}
	example := ""
	if i := strings.Index(desc, "\nExample Usage\n"); i != -1 {
		example = formatHCL(desc[i+15:])
		desc = strings.TrimSpace(desc[:i])
	} else {
		fail(fmt.Sprintf("example usage missing: %s", mdPath))
	}
	descShort := desc
	if i := strings.Index(desc, "\n\n"); i != -1 {
		descShort = strings.TrimSpace(desc[:i])
	}

	out := map[string]string{
		"product":           product,
		"name":              name,
		"dtype":             string(fwList),
		"resource":          resName,
		"cloud_mark":        cloudMark,
		"cloud_title":       cloudTitle,
		"example":           example,
		"description":       desc,
		"description_short": descShort,
		"arguments":         "",
		"attributes":        "",
		"timeouts":          "",
		"import":            "",
	}

	outPath := filepath.Join(docRoot, fwList.outputDir(), fmt.Sprintf("%s.html.markdown", resName))
	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		fail(fmt.Sprintf("mkdir %s failed: %s", filepath.Dir(outPath), err))
	}
	fd, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		fail(fmt.Sprintf("open %s failed: %s", outPath, err))
	}
	defer fd.Close()
	t := template.Must(template.New("fw_list").Parse(docTPL))
	if err := t.Execute(fd, out); err != nil {
		fail(fmt.Sprintf("write %s failed: %s", outPath, err))
	}
	message("[SUCC.]write doc to file success: %s", outPath)
}

// lookupFrameworkMd searches for a .md file matching <name>.md by
// walking framework/*/[type]/<name>.md. The product label hint speeds up
// the lookup but is not required to match.
func lookupFrameworkMd(frameworkRoot string, dtype fwDocType, resName, _ string) string {
	relDir := dtype.sourceMdRelDir()
	if relDir == "" {
		return ""
	}
	var found string
	_ = filepath.Walk(frameworkRoot, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if filepath.Base(filepath.Dir(p)) != relDir {
			return nil
		}
		if info.Name() == resName+".md" {
			found = p
		}
		return nil
	})
	return found
}

// readFrameworkIndex parses tencentcloud/framework/provider.md and
// returns the per-product reference grouping.
func readFrameworkIndex(frameworkRoot string) []frameworkProduct {
	path := filepath.Join(frameworkRoot, "provider.md")
	raw, err := os.ReadFile(path)
	if err != nil {
		fail(fmt.Sprintf("read %s failed: %s", path, err))
	}
	desc := strings.TrimSpace(string(raw))
	pos := strings.Index(desc, "\nResources List\n")
	if pos == -1 {
		fail(fmt.Sprintf("framework provider.md missing 'Resources List' header: %s", path))
	}
	doc := strings.TrimSpace(desc[pos+16:])
	prods, err := getFrameworkIndex(doc)
	if err != nil {
		fail(fmt.Sprintf("parse %s failed: %s", path, err))
	}
	return prods
}

// productOf returns the product label of a typed resource by scanning
// the parsed framework index.
func productOf(prods []frameworkProduct, name string, dtype fwDocType) string {
	for _, p := range prods {
		var bag []string
		switch dtype {
		case fwResource:
			bag = p.Resources
		case fwDataSrc:
			bag = p.DataSources
		case fwFunction:
			bag = p.Functions
		case fwEphemeral:
			bag = p.Ephemerals
		case fwList:
			bag = p.Lists
		case fwAction:
			bag = p.Actions
		}
		for _, n := range bag {
			if n == name {
				return p.Name
			}
		}
	}
	return "Provider Meta"
}

// productOfList is a thin convenience wrapper.
func productOfList(prods []frameworkProduct, name string) string {
	return productOf(prods, name, fwList)
}

// allListNames flattens every list-resource entry in the parsed index.
func allListNames(prods []frameworkProduct) []string {
	var out []string
	for _, p := range prods {
		out = append(out, p.Lists...)
	}
	return out
}

// frameworkResourceTypeName extracts the type name from a framework
// Resource by calling Metadata.
func frameworkResourceTypeName(ctx context.Context, primary frameworkProviderTypeNamer, r resource.Resource) string {
	resp := resource.MetadataResponse{}
	r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: cloudMark}, &resp)
	return resp.TypeName
}

// frameworkDataSourceTypeName extracts the type name of a data source.
func frameworkDataSourceTypeName(ctx context.Context, primary frameworkProviderTypeNamer, ds datasource.DataSource) string {
	resp := datasource.MetadataResponse{}
	ds.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: cloudMark}, &resp)
	return resp.TypeName
}

// frameworkFunctionName extracts the type name of a function. Functions
// expose their name via Metadata too.
func frameworkFunctionName(ctx context.Context, primary frameworkProviderTypeNamer, fn function.Function) string {
	resp := function.MetadataResponse{}
	fn.Metadata(ctx, function.MetadataRequest{}, &resp)
	if !strings.HasPrefix(resp.Name, cloudPrefix) {
		// Function names are by convention emitted without the
		// "tencentcloud_" prefix, but the doc pipeline needs the
		// fully-qualified name to align with the index.
		return cloudPrefix + resp.Name
	}
	return resp.Name
}

// frameworkEphemeralName extracts the type name of an ephemeral resource.
func frameworkEphemeralName(ctx context.Context, primary frameworkProviderTypeNamer, er ephemeral.EphemeralResource) string {
	resp := ephemeral.MetadataResponse{}
	er.Metadata(ctx, ephemeral.MetadataRequest{ProviderTypeName: cloudMark}, &resp)
	return resp.TypeName
}

// frameworkActionName extracts the type name of an action.
func frameworkActionName(ctx context.Context, primary frameworkProviderTypeNamer, a frameworkaction.Action) string {
	resp := frameworkaction.MetadataResponse{}
	a.Metadata(ctx, frameworkaction.MetadataRequest{ProviderTypeName: cloudMark}, &resp)
	return resp.TypeName
}

// frameworkProviderTypeNamer is a forward-compat marker so the helper
// signatures can later carry the SDKv2 provider for cross-reference if
// needed.
type frameworkProviderTypeNamer interface{}

// fail centralises the error path so framework_schema.go and
// framework.go both abort the same way as the sdkv2 generator does.
func fail(msg string) {
	message("[FAIL!]%s", msg)
	os.Exit(1)
}

// keep frameworkprovider import alive — the var assertion makes the
// dependency visible to the linker even though we only call NewProvider
// indirectly through the cloud package.
var _ frameworkprovider.Provider = (*tcfw.Provider)(nil)
