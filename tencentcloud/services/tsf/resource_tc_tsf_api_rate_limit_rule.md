Provides a resource to create a tsf api_rate_limit_rule

Example Usage

```hcl
resource "tencentcloud_tsf_api_rate_limit_rule" "api_rate_limit_rule" {
  api_id = "api-xxxxxx"
  max_qps = 10
  usable_status = "enable"
}
```

Import

tsf api_rate_limit_rule can be imported using the id, e.g.

```
terraform import tencentcloud_tsf_api_rate_limit_rule.api_rate_limit_rule api_rate_limit_rule_id
```