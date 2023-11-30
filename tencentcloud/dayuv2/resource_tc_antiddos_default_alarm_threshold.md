Provides a resource to create a antiddos default alarm threshold

Example Usage

```hcl
resource "tencentcloud_antiddos_default_alarm_threshold" "default_alarm_threshold" {
  default_alarm_config {
	alarm_type = 1
	alarm_threshold = 2000
  }
  instance_type = "bgp"
}
```

Import

antiddos default_alarm_threshold can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_default_alarm_threshold.default_alarm_threshold ${instanceType}
```