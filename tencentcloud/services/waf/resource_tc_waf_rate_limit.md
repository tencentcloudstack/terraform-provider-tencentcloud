Provides a resource to create a WAF rate limit rule

Example Usage

```hcl
resource "tencentcloud_waf_rate_limit" "example" {
  domain         = "example.com"
  name           = "tf-example"
  priority       = 100
  status         = 1
  limit_object   = "Domain"
  limit_strategy = 1

  limit_window {
    second = 10
    minute = 100
  }
}
```

Create with API path rate limiting

```hcl
resource "tencentcloud_waf_rate_limit" "api_example" {
  domain         = "example.com"
  name           = "tf-api-example"
  priority       = 200
  status         = 1
  limit_object   = "API"
  limit_strategy = 1
  object_src     = 0
  order          = 0

  limit_window {
    second = 5
  }

  limit_paths {
    path = "/api/v1/users"
    type = "EXACT"
  }

  limit_method {
    method = "POST"
    type   = "EXACT"
  }

  paths_option {
    path   = "/api/v1/users"
    method = "POST"
  }
}
```

Import

WAF rate limit rule can be imported using the composite id domain#limit_rule_id, e.g.

```
terraform import tencentcloud_waf_rate_limit.example example.com#12345
```
