Provides a resource to create a WAF api sec sensitive privilege rule

Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_privilege_rule" "example" {
  domain         = "www.example.com"
  rule_name      = "tf-example"
  status         = 1
  api_name       = ["/api/user/info"]
  position       = "QUERY"
  parameter_list = ["parameter"]
  option         = 1
  source         = "custom"

  api_name_op {
    op    = "belong"
    value = ["/api/user/info"]
  }
}
```

Import

WAF api sec sensitive privilege rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_privilege_rule.example www.example.com#tf-example
```
