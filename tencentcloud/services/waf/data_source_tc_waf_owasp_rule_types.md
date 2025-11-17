Use this data source to query detailed information of WAF owasp rule types

Example Usage

```hcl
data "tencentcloud_waf_owasp_rule_types" "example" {
  domain = "demo.com"
  filters {
    name        = "RuleId"
    values      = ["10000001"]
    exact_match = true
  }
}
```
