package main

import (
	"strings"
	"testing"
)

func TestSformatHCL(t *testing.T) {
	hcl := `
'hcl
data "tencentcloud_as_scaling_configs" "as_configs" {
configuration_id="asc-oqio4yyj"
result_output_file="my_test_path"
}
'
	`
	hclExp := `'hcl
data "tencentcloud_as_scaling_configs" "as_configs" {
  configuration_id   = "asc-oqio4yyj"
  result_output_file = "my_test_path"
}
'`
	hcl = strings.Replace(hcl, "'", "```", -1)
	hclExp = strings.Replace(hclExp, "'", "```", -1)

	hcl = formatHCL(hcl)
	if hcl != hclExp {
		t.Error("format hcl failed")
	}
}

func TestMformatHCL(t *testing.T) {
	hcl := `
Private Bucket
'hcl
resource "tencentcloud_cos_bucket" "mycos" {
	bucket = "mycos-1258798060"
	acl    = "private"
}
'
Static Website
'hcl
resource "tencentcloud_cos_bucket" "mycos" {
	bucket = "mycos-1258798060"

	website = {
	index_document = "index.html"
	error_document = "error.html"
	}
}
'
	`
	hclExp := `Private Bucket

'hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"
  acl    = "private"
}
'

Static Website

'hcl
resource "tencentcloud_cos_bucket" "mycos" {
  bucket = "mycos-1258798060"

  website = {
    index_document = "index.html"
    error_document = "error.html"
  }
}
'`
	hcl = strings.Replace(hcl, "'", "```", -1)
	hclExp = strings.Replace(hclExp, "'", "```", -1)

	hcl = formatHCL(hcl)
	if hcl != hclExp {
		t.Error("format hcl failed")
	}
}

func TestContainsBigSymbol(t *testing.T) {
	cases := "中国人繁體字ａｂｃＡＢＣ～！＠＃￥％…＆（）—＋｛｝｜：“”《》？１２３４５６７８９０－＝【】＼；‘’，。、　"
	for _, c := range cases {
		if ContainsBigSymbol(string(c)) == "" {
			t.Log(c)
			t.Errorf("Expected %s to be Chinese symbol", string(c))
		}
	}

	cases = "abcABC~!@#$%^&*()_+{}|:\"<>?`1234567890-=[]\\;',./ \t\r\n"
	for _, c := range cases {
		if ContainsBigSymbol(string(c)) != "" {
			t.Errorf("Expected %s not to be Chinese symbol", string(c))
		}
	}
}
