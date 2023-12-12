Provides a resource to create a css enable_optimal_switching

~> **NOTE:** This resource is only valid when the push stream. When the push stream ends, it will be deleted.

Example Usage

```hcl
resource "tencentcloud_css_enable_optimal_switching" "enable_optimal_switching" {
  stream_name     = "1308919341_test"
  enable_switch   = 1
  host_group_name = "test-group"
}
```

Import

css domain can be imported using the id, e.g.

```
terraform import tencentcloud_css_enable_optimal_switching.enable_optimal_switching streamName
```