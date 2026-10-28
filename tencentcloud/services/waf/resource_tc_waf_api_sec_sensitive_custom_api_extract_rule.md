Provides a resource to create a WAF api sec sensitive custom api extract rule

Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule" "example" {
  domain    = "www.example.com"
  rule_name = "tf-example"
  status    = 1
  api_name  = "/api/login"
  methods   = ["GET", "POST"]
  regex     = "/api/.*"
}
```

Import

WAF api sec sensitive custom api extract rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_custom_api_extract_rule.example www.example.com#tf-example
```
