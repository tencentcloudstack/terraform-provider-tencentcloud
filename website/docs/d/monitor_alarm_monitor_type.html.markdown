---
subcategory: "Cloud Monitor(Monitor)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_monitor_alarm_monitor_type"
sidebar_current: "docs-tencentcloud-datasource-monitor_alarm_monitor_type"
description: |-
  Use this data source to query detailed information of monitor alarm_monitor_type
---

# tencentcloud_monitor_alarm_monitor_type

Use this data source to query detailed information of monitor alarm_monitor_type

## Example Usage

```hcl
data "tencentcloud_monitor_alarm_monitor_type" "alarm_monitor_type" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `monitor_type_infos` - Monitoring type details.
  * `id` - Monitoring type ID.
  * `name` - Monitoring type.
  * `sort_id` - Sort order.
* `monitor_types` - Monitoring type, cloud product monitoring is MT_ QCE.


