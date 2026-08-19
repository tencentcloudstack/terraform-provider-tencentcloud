Provides a resource to create a WAF api sec sensitive custom api exclude rule

Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule" "example" {
  domain     = "www.example.com"
  rule_name  = "tf-example"
  status     = 1
  match_type = "regex"
  content    = "/static"
}
```

Import

WAF api sec sensitive custom api exclude rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_custom_api_exclude_rule.example www.example.com#tf-example
```
