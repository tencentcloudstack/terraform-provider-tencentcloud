Provides a resource to create a WAF api sec sensitive custom rule

Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_custom_rule" "example" {
  domain      = "www.example.com"
  rule_name   = "tf-example"
  status      = 1
  position    = ["headers"]
  match_key   = "key_match"
  match_value = ["admin", "cookie"]
  level       = "100"
  match_cond  = ["and"]
  is_pan      = 1
}
```

Import

WAF api sec sensitive custom rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_custom_rule.example www.example.com#tf-example
```
