---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_time_window"
sidebar_current: "docs-tencentcloud-resource-mysql_time_window"
description: |-
  Provides a resource to create a mysql time_window
---

# tencentcloud_mysql_time_window

Provides a resource to create a mysql time_window

## Example Usage

```hcl
resource "tencentcloud_mysql_time_window" "time_window" {
  instance_id    = "cdb-lw71b6ar"
  max_delay_time = 10
  time_ranges = [
    "01:00-02:01"
  ]
  weekdays = [
    "friday",
    "monday",
    "saturday",
    "thursday",
    "tuesday",
    "wednesday",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID in the format of cdb-c1nl9rpv or cdbro-c1nl9rpv. It is the same as the instance ID displayed on the TencentDB Console page.
* `time_ranges` - (Required, Set: [`String`]) Time period available for maintenance after modification in the format of 10:00-12:00. Each period lasts from half an hour to three hours, with the start time and end time aligned by half-hour. Up to two time periods can be set. Start and end time range: [00:00, 24:00].
* `max_delay_time` - (Optional, Int) Data delay threshold. It takes effect only for source instance and disaster recovery instance. Default value: 10.
* `weekdays` - (Optional, Set: [`String`]) Specifies for which day to modify the time period. Value range: Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday. If it is not specified or is left blank, the time period will be modified for every day by default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mysql time_window can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_time_window.time_window instanceId
```

