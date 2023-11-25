---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_default_alarm_threshold"
sidebar_current: "docs-tencentcloud-resource-antiddos_default_alarm_threshold"
description: |-
  Provides a resource to create a antiddos default alarm threshold
---

# tencentcloud_antiddos_default_alarm_threshold

Provides a resource to create a antiddos default alarm threshold

## Example Usage

```hcl
resource "tencentcloud_antiddos_default_alarm_threshold" "default_alarm_threshold" {
  default_alarm_config {
    alarm_type      = 1
    alarm_threshold = 2000
  }
  instance_type = "bgp"
}
```

## Argument Reference

The following arguments are supported:

* `default_alarm_config` - (Required, List) Alarm threshold configuration.
* `instance_type` - (Required, String, ForceNew) Product type, value [bgp (represents advanced defense package product) bgpip (represents advanced defense IP product)].

The `default_alarm_config` object supports the following:

* `alarm_threshold` - (Optional, Int) Alarm threshold, in Mbps, with a value of&gt;=0; When used as an input parameter, setting 0 will delete the alarm threshold configuration;.
* `alarm_type` - (Optional, Int) Alarm threshold type, value [1 (incoming traffic alarm threshold) 2 (attack cleaning traffic alarm threshold)].

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

antiddos default_alarm_threshold can be imported using the id, e.g.

```
terraform import tencentcloud_antiddos_default_alarm_threshold.default_alarm_threshold ${instanceType}
```

