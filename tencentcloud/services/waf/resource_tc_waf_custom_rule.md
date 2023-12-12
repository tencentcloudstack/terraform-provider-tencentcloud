Provides a resource to create a waf custom_rule

Example Usage

```hcl
resource "tencentcloud_waf_custom_rule" "example" {
  name        = "tf-example"
  sort_id     = "50"
  redirect    = "/"
  expire_time = "0"

  strategies {
    field        = "IP"
    compare_func = "ipmatch"
    content      = "2.2.2.2"
    arg          = ""
  }

  status      = "1"
  domain      = "test.com"
  action_type = "1"
}
```

Import

waf custom_rule can be imported using the id, e.g.

```
terraform import tencentcloud_waf_custom_rule.example test.com#1100310609
```