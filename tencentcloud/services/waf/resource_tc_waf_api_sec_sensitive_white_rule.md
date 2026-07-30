Provides a resource to create a WAF api sec sensitive white rule

Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_white_rule" "example" {
  domain      = "www.example.com"
  rule_name   = "tf-example"
  status      = 1
  white_mode  = 1
  description = "tf example white rule"

  api_name_op {
    op    = "belong"
    value = ["/api/user/info"]
  }
}
```

Import

WAF api sec sensitive white rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_white_rule.example www.example.com#tf-example
```
