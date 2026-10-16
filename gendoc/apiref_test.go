package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderAPIRefs(t *testing.T) {
	refs := map[string][]apiRef{
		"Delete": {{Product: "vpc", Action: "DeleteVpc", URL: "https://example.com/delete"}},
		"Create": {{Product: "vpc", Action: "CreateVpc", URL: "https://example.com/create"}},
		"Read":   {{Product: "vpc", Action: "DescribeVpcs", URL: ""}},
	}
	got := renderAPIRefs(refs)
	want := "This resource interacts with the following TencentCloud APIs:\n\n" +
		"* Create: [CreateVpc](https://example.com/create)\n" +
		"* Read: `DescribeVpcs`\n" +
		"* Delete: [DeleteVpc](https://example.com/delete)"
	if got != want {
		t.Errorf("renderAPIRefs mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderAPIRefsEmpty(t *testing.T) {
	if got := renderAPIRefs(nil); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestIsAPIActionName(t *testing.T) {
	cases := []struct {
		in   string
		want bool
	}{
		{"CreateVpc", true},
		{"CreateVpcWithContext", true}, // accepted, suffix stripped later
		{"callWithRetry", false},
		{"AwaitTask", false},
		{"WithHttpTransport", false},
		{"unexported", false},
		{"WithContext", false}, // pure suffix, stripping yields empty
	}
	for _, c := range cases {
		if got := isAPIActionName(c.in); got != c.want {
			t.Errorf("isAPIActionName(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestStripContextSuffix(t *testing.T) {
	if got := stripContextSuffix("CreateVpcWithContext"); got != "CreateVpc" {
		t.Errorf("stripContextSuffix = %q", got)
	}
	if got := stripContextSuffix("CreateVpc"); got != "CreateVpc" {
		t.Errorf("stripContextSuffix = %q", got)
	}
}

func TestRecvTypeName(t *testing.T) {
	if got := recvTypeName(nil); got != "" {
		t.Errorf("recvTypeName(nil) = %q", got)
	}
}

func TestRenderResourceAPISectionManualOverride(t *testing.T) {
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Skip("no repo root")
	}
	if _, err := os.Stat(filepath.Join(repoRoot, "tencentcloud", "services")); err != nil {
		t.Skip("not run from gendoc dir")
	}
	manual := "* Create: MyCustomApi"
	got := renderResourceAPISection(nil, repoRoot, "tencentcloud_vpc", nil, manual)
	if got != manual {
		t.Errorf("manual override not honored: got %q", got)
	}
}

func TestRenderAPIRefsNamelessAndLinkless(t *testing.T) {
	// Actions render as "Verb: [Action](url)" only — no descriptions. An
	// empty URL degrades to a bare code-quoted action name.
	refs := map[string][]apiRef{
		"Create": {
			{Product: "vpc", Action: "CreateVpc", URL: "https://example.com/create"},
			{Product: "vpc", Action: "NoUrl", URL: ""},
		},
	}
	got := renderAPIRefsFor("resource", refs)
	want := "This resource interacts with the following TencentCloud APIs:\n\n" +
		"* Create: [CreateVpc](https://example.com/create)\n" +
		"* Create: `NoUrl`"
	if got != want {
		t.Errorf("render mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
	// Descriptions were removed from the section; none may leak back in.
	if strings.Contains(got, " - ") {
		t.Error("description text leaked into the rendered section")
	}
}

func TestResolveAPISearchURLFallback(t *testing.T) {
	// An action without a catalogue page falls back to the international
	// site's search page for the action name.
	a := &apiAnalyzer{
		clientProd: map[string]string{"accesskey": "cam"},
		catalog:    &catalog{urlIdx: map[string]string{}},
	}
	ref := a.resolveAPI(apiCallKey{client: "accesskey", action: "CreateAccessKey"})
	if ref == nil {
		t.Fatal("expected a ref")
	}
	want := "https://www.tencentcloud.com/search/?wd=CreateAccessKey&type=document&lang=en"
	if ref.URL != want {
		t.Errorf("fallback URL = %q, want %q", ref.URL, want)
	}
}

func TestResolveAPICatalogHit(t *testing.T) {
	// An action present in the international catalogue resolves to its
	// docs page URL.
	a := &apiAnalyzer{
		clientProd: map[string]string{"accesskey": "cam"},
		catalog: &catalog{urlIdx: map[string]string{
			"cam.CreateAccessKey": "https://www.tencentcloud.com/document/api/1302/47414",
		}},
	}
	ref := a.resolveAPI(apiCallKey{client: "accesskey", action: "CreateAccessKey"})
	if ref == nil || ref.URL != "https://www.tencentcloud.com/document/api/1302/47414" {
		t.Fatalf("catalog hit failed: %+v", ref)
	}
}

func TestRenderAPIRefsSkipsEmptyVerbLists(t *testing.T) {
	// A map carrying only empty lists must render nothing, not a bare
	// lead-in sentence (P2 regression guard).
	refs := map[string][]apiRef{"Read": nil}
	if got := renderAPIRefsFor("resource", refs); got != "" {
		t.Errorf("expected empty, got %q", got)
	}
}

func TestBalanceFences(t *testing.T) {
	// Real upstream shape: the Example Usage block is unclosed, so the
	// Import fence markers pair with it; the lone stray ``` in the middle
	// must be the one dropped, and the trailing fence must survive.
	in := "prose\n\n```hcl\nresource \"x\" {}\n\nImport\n\ne.g.\n\n```\nterraform import x y\n```\n\n```"
	got := balanceFences(in)
	if strings.Count(got, "```")%2 != 0 {
		t.Errorf("still unbalanced: %q", got)
	}
	if !strings.Contains(got, "terraform import x y") {
		t.Errorf("content lost: %q", got)
	}
	// A language-tagged opener (```hcl) must never be the one escaped.
	if !strings.Contains(got, "```hcl") {
		t.Errorf("lost the language-tagged opener: %q", got)
	}
	// The stray marker is escaped to literal text, not deleted.
	if !strings.Contains(got, "`"+"`"+"`") {
		t.Errorf("stray fence not escaped: %q", got)
	}
	// Balanced input passes through unchanged.
	ok := "```\ncode\n```"
	if got := balanceFences(ok); got != ok {
		t.Errorf("balanced input modified: %q", got)
	}
}
