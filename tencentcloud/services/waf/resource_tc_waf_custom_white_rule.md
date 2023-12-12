Provides a resource to create a waf custom_white_rule

Example Usage

```hcl
resource "tencentcloud_waf_custom_white_rule" "example" {
  name        = "tf-example"
  sort_id     = "30"
  expire_time = "0"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "1.1.1.1"
    arg          = ""
  }

  status = "1"
  domain = "test.com"
  bypass = "geoip,cc,owasp"
}
```

Import

waf custom_white_rule can be imported using the id, e.g.

```
terraform import tencentcloud_waf_custom_white_rule.example test.com#1100310837
```