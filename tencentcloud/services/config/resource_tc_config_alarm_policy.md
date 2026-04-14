Provides a resource to create a Config alarm policy.

Example Usage

```hcl
resource "tencentcloud_config_alarm_policy" "example" {
  name                   = "tf-example"
  event_scope            = [1]
  risk_level             = [1, 2]
  notice_time            = "09:30:00~23:30:00"
  notification_mechanism = "实时发送"
  status                 = 1
  notice_period          = [1, 2, 3, 4, 5]
  description            = "tf example alarm policy"
}
```

Import

Config alarm policy can be imported using the alarmPolicyId, e.g.

```
terraform import tencentcloud_config_alarm_policy.example 123
```
