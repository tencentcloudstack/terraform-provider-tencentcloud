Provides a resource to create a WAF api sec sensitive custom event rule

Example Usage

```hcl
resource "tencentcloud_waf_api_sec_sensitive_custom_event_rule" "example" {
  domain        = "www.example.com"
  rule_name     = "tf-example"
  status        = 1
  description   = "tf example custom event rule"
  req_frequency = [10, 1]
  risk_level    = "100"
  source        = "custom"

  api_name_op {
    op    = "belong"
    value = ["/api/login"]

    api_name_method {
      api_name = "/api/login"
      method   = "POST"
    }
  }

  match_rule_list {
    key     = "get_key"
    operate = "exist"
    value   = ["admin", "root"]
  }

  stat_rule_list {
    key     = "status"
    operate = "num_gt"
    value   = ["50"]
    name    = "200"
  }
}
```

Import

WAF api sec sensitive custom event rule can be imported using the domain#ruleName, e.g.

```
terraform import tencentcloud_waf_api_sec_sensitive_custom_event_rule.example www.example.com#tf-example
```
