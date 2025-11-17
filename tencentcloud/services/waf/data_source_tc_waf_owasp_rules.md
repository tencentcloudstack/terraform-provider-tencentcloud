Use this data source to query detailed information of WAF owasp rules

Example Usage

```hcl
data "tencentcloud_waf_owasp_rules" "example" {
  domain = "example.qcloud.com"
  by     = "RuleId"
  order  = "desc"
  filters {
    name        = "RuleId"
    values      = ["106251141"]
    exact_match = true
  }
}
```
