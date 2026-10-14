package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestRenderAPIRefs(t *testing.T) {
	refs := map[string][]apiRef{
		"Delete": {{Product: "vpc", Action: "DeleteVpc", Name: "删除VPC", URL: "https://example.com/delete"}},
		"Create": {{Product: "vpc", Action: "CreateVpc", Name: "创建VPC", URL: "https://example.com/create"}},
		"Read":   {{Product: "vpc", Action: "DescribeVpcs", Name: "查询VPC列表", URL: ""}},
	}
	got := renderAPIRefs(refs)
	want := "This resource interacts with the following TencentCloud APIs:\n\n" +
		"* Create: [CreateVpc](https://example.com/create) - 创建VPC\n" +
		"* Read: `DescribeVpcs`\n" +
		"* Delete: [DeleteVpc](https://example.com/delete) - 删除VPC"
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

func TestRenderAPIRefsPrefersEnglish(t *testing.T) {
	refs := map[string][]apiRef{
		"Create": {
			{Product: "vpc", Action: "CreateVpc", Name: "Creates a VPC instance.", URL: "https://example.com/create"},
			{Product: "vpc", Action: "NoEnDesc", Name: "", URL: "https://example.com/noen"},
			{Product: "vpc", Action: "NoUrl", Name: "Has name but no url.", URL: ""},
		},
	}
	got := renderAPIRefsFor("resource", refs)
	want := "This resource interacts with the following TencentCloud APIs:\n\n" +
		"* Create: [CreateVpc](https://example.com/create) - Creates a VPC instance.\n" +
		"* Create: [NoEnDesc](https://example.com/noen)\n" +
		"* Create: `NoUrl`"
	if got != want {
		t.Errorf("render mismatch:\ngot:\n%s\nwant:\n%s", got, want)
	}
	// Generated docs must never contain the Chinese catalog description.
	if strings.Contains(got, "创建") {
		t.Error("Chinese text leaked into the rendered section")
	}
}

func TestExtractENFromClient(t *testing.T) {
	src := `package vpc

// CreateVpc
// This API is used to create a VPC instance.
//
// error code that may be returned:
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) CreateVpc(request *CreateVpcRequest) (response *CreateVpcResponse, err error) {
	return nil, nil
}

// CreateVpc
// This API is used to create a VPC instance.
//
// error code that may be returned:
//  INVALIDPARAMETER = "InvalidParameter"
func (c *Client) CreateVpcWithContext(ctx context.Context, request *CreateVpcRequest) (response *CreateVpcResponse, err error) {
	return nil, nil
}
`
	got := extractENFromClient(src)
	if len(got) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(got))
	}
	if got["CreateVpc"] != "This API is used to create a VPC instance." {
		t.Errorf("desc = %q", got["CreateVpc"])
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

func TestCleanDescTruncation(t *testing.T) {
	// A description whose cut point lands inside a markdown link must not
	// produce broken markdown (P3: unbalanced parentheses).
	long := "This API is used to do something. See [RunInstances](https://intl.cloud.tencent.com/document/product/213/33237) for more information about the instance launch flow and the parameters that it accepts."
	got := cleanDesc(long)
	if strings.Count(got, "(") != strings.Count(got, ")") {
		t.Errorf("unbalanced parens: %q", got)
	}
	if strings.Count(got, "[") != strings.Count(got, "]") {
		t.Errorf("unbalanced brackets: %q", got)
	}
	if utf8.RuneCountInString(got) > maxDescLen {
		t.Errorf("too long: %d runes", utf8.RuneCountInString(got))
	}
	// Multi-byte characters must not be split mid-rune.
	multi := strings.Repeat("描述文字", 200)
	got2 := cleanDesc(multi)
	if strings.HasSuffix(got2, "\ufffd") {
		t.Error("split a multi-byte rune")
	}
}

func TestCleanDescShort(t *testing.T) {
	s := "Short description."
	if got := cleanDesc(s); got != s {
		t.Errorf("short desc changed: %q", got)
	}
}

func TestCleanDescTeoCase(t *testing.T) {
	// Reproduce the teo.DescribeSecurityIPGroupInfo leftover: the sentence
	// boundary cut must not leave an unclosed "[" in the output.
	long := "The API is deprecated and will be discontinued on June 30, 2024. Please use the API [DescribeSecurityIPGroups](https://www.tencentcloud.com/document/product/1552/81890) instead to manage security IP groups for the web/mini-program business."
	got := cleanDesc(long)
	t.Logf("got: %q", got)
	if strings.Count(got, "[") != strings.Count(got, "]") {
		t.Errorf("unbalanced brackets: %q", got)
	}
}

func TestCleanDescURLPeriod(t *testing.T) {
	// The sentence's link URL exceeds the cap: the last "." inside the cut
	// is a dot inside the URL, and the sentence boundary fallback must
	// still avoid emitting an unclosed "[".
	long := "This API is on an earlier version and will be discontinued in a future release. If you want to call it, please switch to the latest version. For details, see [ModifyL4ProxyRulesStatus](https://www.tencentcloud.com/document/product/1552/81891) which replaces this deprecated interface for all web and mini-program businesses worldwide."
	got := cleanDesc(long)
	t.Logf("got: %q", got)
	if strings.Count(got, "[") != strings.Count(got, "]") {
		t.Errorf("unbalanced brackets: %q", got)
	}
}
