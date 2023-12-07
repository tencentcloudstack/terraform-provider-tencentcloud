Provides a resource to create a waf anti_info_leak

Example Usage

```hcl
resource "tencentcloud_waf_anti_info_leak" "example" {
  domain      = "tf.example.com"
  name        = "tf_example"
  action_type = 0
  strategies {
    field   = "information"
    content = "phone"
  }
  uri    = "/anti_info_leak_url"
  status = 1
}
```

Import

waf anti_info_leak can be imported using the id, e.g.

```
terraform import tencentcloud_waf_anti_info_leak.example 3100077499#tf.example.com
```