package main

import (
	"os"
	"path/filepath"
	"testing"
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
