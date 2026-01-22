Provides a resource to create a WAF owasp rule type config

Example Usage

```hcl
resource "tencentcloud_waf_owasp_rule_type_config" "example" {
  domain           = "demo.com"
  type_id          = "30000000"
  rule_type_status = 1
  rule_type_action = 1
  rule_type_level  = 200
}
```

Import

WAF owasp rule type config can be imported using the domain#typeId, e.g.

```
terraform import tencentcloud_waf_owasp_rule_type_config.example demo.com#30000000
```
