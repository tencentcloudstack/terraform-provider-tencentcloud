---
subcategory: "Anti-DDoS(DayuV2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_antiddos_overview_ddos_event_list"
sidebar_current: "docs-tencentcloud-datasource-antiddos_overview_ddos_event_list"
description: |-
  Use this data source to query detailed information of antiddos overview_ddos_event_list
---

# tencentcloud_antiddos_overview_ddos_event_list

Use this data source to query detailed information of antiddos overview_ddos_event_list

## Example Usage

```hcl
data "tencentcloud_antiddos_overview_ddos_event_list" "overview_ddos_event_list" {
  start_time    = "2023-11-20 00:00:00"
  end_time      = "2023-11-21 00:00:00"
  attack_status = "end"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) EndTime.
* `start_time` - (Required, String) StartTime.
* `attack_status` - (Optional, String) filter event by attack status, start: attacking; end: attack end.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `event_list` - EventList.
  * `attack_status` - Attack status, 0: Under attack; 1: End of attack.
  * `attack_type` - AttackType.
  * `business` - Dayu sub product code (bgpip represents advanced defense IP; net represents professional version of advanced defense IP).
  * `end_time` - EndTime.
  * `id` - event id.
  * `instance_id` - InstanceId.
  * `instance_name` - InstanceId.
  * `mbps` - Attack traffic, unit Mbps.
  * `pps` - unit Mbps.
  * `start_time` - StartTime.
  * `vip` - ip.


