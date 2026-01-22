Provides a resource to create a WAF owasp rule status config

Example Usage

```hcl
resource "tencentcloud_waf_owasp_rule_status_config" "example" {
  domain      = "demo.com"
  rule_id     = "106251141"
  rule_status = 1
}
```

Import

WAF owasp rule status config can be imported using the domain#ruleId, e.g.

```
terraform import tencentcloud_waf_owasp_rule_status_config.example demo.com#106251141
```
