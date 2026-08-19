// Package main — framework.go is the framework-stack counterpart of the
// sdkv2 documentation generator. It enumerates the six framework
// reference types (Resource / DataSource / Function / EphemeralResource /
// ListResource / Action), reads their per-type markdown files placed
// next to the Go source under tencentcloud/services/<product>/, and
// renders them into website/docs/<dir>/<resource>.html.markdown.
//
// Output dirs (single-letter where possible, mirroring the SDKv2 d/r
// shorthand):
//
//	resource           -> website/docs/r/
//	datasource         -> website/docs/d/
//	function           -> website/docs/f/
//	ephemeral resource -> website/docs/e/
//	list resource      -> website/docs/l/
//	action             -> website/docs/a/
//
// The framework references are listed inside the unified index file
// tencentcloud/provider.md alongside SDKv2 references. GetIndex (in
// index.go) handles the six section headers; this file only consumes
// the parsed []Product slice via genFrameworkDocs.
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
// The single-letter shorthands mirror the SDKv2 d/r convention.
func (t fwDocType) outputDir() string {
	switch t {
	case fwResource:
		return "r"
	case fwDataSrc:
		return "d"
	case fwFunction:
		return "f"
	case fwEphemeral:
		return "e"
	case fwList:
		return "l"
	case fwAction:
		return "a"
	}
	return ""
}

// sourceMdRelDir is no longer used as a directory name (framework
// references now live alongside SDKv2 code under
// tencentcloud/services/<product>/). It is retained as documentation of
// the legacy directory naming for reference only.
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

// mdFilePrefix returns the file-name prefix used by the new
// services/<product>/ layout. Combined with the resource short name and
// an optional product segment, it forms the full <prefix><resName>.md or
// <prefix><product>_<resName>.md filename.
func (t fwDocType) mdFilePrefix() string {
	switch t {
	case fwResource:
		return "resource_tc_"
	case fwDataSrc:
		return "data_source_tc_"
	case fwFunction:
		return "function_tc_"
	case fwEphemeral:
		return "ephemeral_tc_"
	case fwList:
		return "list_tc_"
	case fwAction:
		return "action_tc_"
	}
	return ""
}

// genFrameworkDocs is the framework-side equivalent of the sdkv2 main
// loop. It is invoked from main() with the unified product list parsed
// from tencentcloud/provider.md.
func genFrameworkDocs(repoRoot string, products []Product) {
	servicesRoot := filepath.Join(repoRoot, "services")

	// 1) Build framework provider so we can drive 6 factory aggregators.
	primary := cloud.Provider()
	prov := tcfw.NewProvider(primary)
	ctx := context.Background()

	// 2) Render each reference type.
	for _, fac := range prov.Resources(ctx) {
		r := fac()
		name := frameworkResourceTypeName(ctx, primary, r)
		genFrameworkDoc(servicesRoot, fwResource, name, productOf(products, name, fwResource), func(out *map[string]fwAttrSpec) string {
			resp := resource.SchemaResponse{}
			r.Schema(ctx, resource.SchemaRequest{}, &resp)
			*out = flattenResourceSchema(resp.Schema)
			return strings.TrimSpace(resp.Schema.Description)
		})
	}
	for _, fac := range prov.DataSources(ctx) {
		ds := fac()
		name := frameworkDataSourceTypeName(ctx, primary, ds)
		genFrameworkDoc(servicesRoot, fwDataSrc, name, productOf(products, name, fwDataSrc), func(out *map[string]fwAttrSpec) string {
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
			genFrameworkDoc(servicesRoot, fwFunction, name, productOf(products, name, fwFunction), func(out *map[string]fwAttrSpec) string {
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
			genFrameworkDoc(servicesRoot, fwEphemeral, name, productOf(products, name, fwEphemeral), func(out *map[string]fwAttrSpec) string {
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
	// under services/<product>/ are however parsed via the unified index
	// so that, once the real list reference is wired up, the
	// documentation pipeline keeps working with no changes.
	for _, name := range allListNames(products) {
		genFrameworkListPlaceholder(servicesRoot, name, productOf(products, name, fwList))
	}
	_ = frameworklist.ListResource(nil) // keep the import alive for future use
	if pa, ok := prov.(frameworkprovider.ProviderWithActions); ok {
		for _, fac := range pa.Actions(ctx) {
			a := fac()
			name := frameworkActionName(ctx, primary, a)
			genFrameworkDoc(servicesRoot, fwAction, name, productOf(products, name, fwAction), func(out *map[string]fwAttrSpec) string {
				resp := frameworkaction.SchemaResponse{}
				a.Schema(ctx, frameworkaction.SchemaRequest{}, &resp)
				*out = flattenActionSchema(resp.Schema)
				return strings.TrimSpace(resp.Schema.Description)
			})
		}
	}
}

// schemaExtractor is the closure passed to genFrameworkDoc to delay the
// per-type-specific schema extraction until product/path lookups are
// done. It populates *out and returns the schema-level description.
type schemaExtractor func(out *map[string]fwAttrSpec) string

// genFrameworkDoc renders a single .md file for any framework reference
// type. It mirrors the sdkv2 genDoc function but reads the inline .md
// from services/<product>/<dtype>_tc_[<product>_]<resource>.md and writes
// the rendered markdown to website/docs/<outputDir>/<name>.html.markdown.
func genFrameworkDoc(servicesRoot string, dtype fwDocType, name, product string, extract schemaExtractor) {
	if name == "" {
		fail(fmt.Sprintf("framework %s has empty type name", dtype))
	}
	if !strings.HasPrefix(name, cloudPrefix) {
		fail(fmt.Sprintf("framework %s %q must start with %q", dtype, name, cloudPrefix))
	}
	resName := name[len(cloudPrefix):]

	// Locate the colocated .md file under services/. The product label
	// from framework provider.md is purely a UI label; the on-disk file
	// is matched by prefix (`<dtype>_tc_`) plus a suffix that may or may
	// not include a product segment (e.g. both `action_tc_reboot_instance.md`
	// and `action_tc_cvm_reboot_instance.md` are accepted).
	mdPath := lookupFrameworkMd(servicesRoot, dtype, resName, product)
	if mdPath == "" {
		fail(fmt.Sprintf("framework %s %q is missing its .md file under services/<product>/ (expected %s[<product>_]%s.md)", dtype, name, dtype.mdFilePrefix(), resName))
	}

	message("[START]get description from file: %s\n", relMdPath(mdPath))

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
func genFrameworkListPlaceholder(servicesRoot, name, product string) {
	resName := strings.TrimPrefix(name, cloudPrefix)
	mdPath := lookupFrameworkMd(servicesRoot, fwList, resName, product)
	if mdPath == "" {
		fail(fmt.Sprintf("framework list %q is missing its .md file under services/<product>/ (expected list_tc_[<product>_]%s.md)", name, resName))
	}
	message("[START]get description from file: %s\n", relMdPath(mdPath))
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

// relMdPath converts an absolute markdown path under tencentcloud/services/
// into a path relative to tencentcloud/ (e.g. services/ssm/ephemeral_tc_ssm_secret_version.md),
// matching the format SDKv2's [START] log uses.
func relMdPath(absPath string) string {
	const seg = string(filepath.Separator) + "services" + string(filepath.Separator)
	if i := strings.Index(absPath, seg); i >= 0 {
		return strings.TrimPrefix(absPath[i+1:], "")
	}
	return absPath
}

// lookupFrameworkMd searches for a .md file matching the new
// services/<product>/ naming convention:
//
//	services/<product>/<prefix><resName>.md
//	services/<product>/<prefix><productSegment>_<resName>.md
//
// where <prefix> is dtype.mdFilePrefix() (e.g. "action_tc_"). The
// product label hint is currently unused; we walk the tree once and
// match by file name. The walk stops at the first hit.
func lookupFrameworkMd(servicesRoot string, dtype fwDocType, resName, _ string) string {
	prefix := dtype.mdFilePrefix()
	if prefix == "" {
		return ""
	}
	// Two acceptable filename shapes:
	//   1. <prefix><resName>.md            (e.g. resource_tc_local_note.md)
	//   2. <prefix>*_<resName>.md          (e.g. action_tc_cvm_reboot_instance.md)
	exact := prefix + resName + ".md"
	suffix := "_" + resName + ".md"

	var found string
	_ = filepath.Walk(servicesRoot, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		fn := info.Name()
		if fn == exact {
			found = p
			return filepath.SkipDir
		}
		if strings.HasPrefix(fn, prefix) && strings.HasSuffix(fn, suffix) {
			// Avoid matching unrelated names where the suffix coincides
			// (e.g. "*_reboot_instance.md" should not match "instance.md").
			if found == "" {
				found = p
			}
		}
		return nil
	})
	return found
}

// productOf returns the product label of a typed framework reference
// by scanning the unified Product list parsed from provider.md. If no
// matching entry is found, "Provider Meta" is returned as the fallback
// product label — mirroring the original framework-only behaviour.
func productOf(prods []Product, name string, dtype fwDocType) string {
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

// allListNames flattens every list-resource entry in the parsed index.
func allListNames(prods []Product) []string {
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
