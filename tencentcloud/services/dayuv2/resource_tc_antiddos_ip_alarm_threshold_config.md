Provides a resource to create a antiddos ip_alarm_threshold_config

Example Usage

```hcl
resource "tencentcloud_antiddos_ip_alarm_threshold_config" "ip_alarm_threshold_config" {
  alarm_type = 1
  alarm_threshold = 2
  instance_ip = "xxx.xxx.xxx.xxx"
  instance_id = "bgp-xxxxxx"
}
```

Import

antiddos ip_alarm_threshold_config can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_ip_alarm_threshold_config.ip_alarm_threshold_config ${instanceId}#${instanceIp}#${alarmType}
```