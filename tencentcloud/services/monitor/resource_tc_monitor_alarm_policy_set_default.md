Provides a resource to create a monitor policy_set_default

Example Usage

```hcl
resource "tencentcloud_monitor_alarm_policy_set_default" "policy_set_default" {
  module    = "monitor"
  policy_id = "policy-u4iykjkt"
}
```