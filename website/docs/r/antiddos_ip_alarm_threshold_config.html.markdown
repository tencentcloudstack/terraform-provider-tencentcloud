---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_ip_alarm_threshold_config"
sidebar_current: "docs-tencentcloud-resource-antiddos_ip_alarm_threshold_config"
description: |-
  Provides a resource to create a antiddos ip_alarm_threshold_config
---

# tencentcloud_antiddos_ip_alarm_threshold_config

Provides a resource to create a antiddos ip_alarm_threshold_config

## Example Usage

```hcl
resource "tencentcloud_antiddos_ip_alarm_threshold_config" "ip_alarm_threshold_config" {
  alarm_type      = 1
  alarm_threshold = 2
  instance_ip     = "xxx.xxx.xxx.xxx"
  instance_id     = "bgp-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `alarm_threshold` - (Required, Int) Alarm threshold, in Mbps, with a value of&gt;=0; When used as an input parameter, setting 0 will delete the alarm threshold configuration;.
* `alarm_type` - (Required, Int) Alarm threshold type, value [1 (incoming traffic alarm threshold) 2 (attack cleaning traffic alarm threshold)].
* `instance_id` - (Required, String) Instance id.
* `instance_ip` - (Required, String) Instance ip.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

antiddos ip_alarm_threshold_config can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_ip_alarm_threshold_config.ip_alarm_threshold_config ${instanceId}#${instanceIp}#${alarmType}
```

