Provides a resource to create a WAF rate limit rule

Example Usage

Create with API path rate limiting

```hcl
resource "tencentcloud_waf_rate_limit" "example" {
  domain         = "example.com"
  name           = "tf-example"
  priority       = 10
  status         = 0
  limit_strategy = 0
  limit_object   = "API"
  block_page     = 209057

  get_params_name {
    content = "get"
    func    = "IN"
  }

  limit_headers {
    key   = "myKey"
    type  = "IN"
    value = "myValue"
  }

  limit_paths {
    path = "/url"
    type = "EXACT"
  }

  limit_window {
    second = 0
    minute = 10
    hour   = 0
  }
}
```

Import

WAF rate limit rule can be imported using the composite id domain#limit_rule_id, e.g.

```
terraform import tencentcloud_waf_rate_limit.example example.com#4000077639
```
