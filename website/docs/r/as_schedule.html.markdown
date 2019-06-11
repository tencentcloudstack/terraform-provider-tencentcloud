---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_schedule"
sidebar_current: "docs-tencentcloud-resource-as_schedule"
description: |-
  Provides a resource for an AS (Auto scaling) schedule.
---

# tencentcloud_as_schedule

Provides a resource for an AS (Auto scaling) schedule.

## Example Usage

```hcl
resource "tencentcloud_as_schedule" "schedule" {
	scaling_group_id = "sg-12af45"
	schedule_action_name = "tf-as-schedule"
	max_size = 10
	min_size = 0
	desired_capacity = 0
	start_time = "2019-01-01T00:00:00+08:00"
	end_time = "2019-12-01T00:00:00+08:00"
	recurrence = "0 0 * * *"
}
```

## Argument Reference

The following arguments are supported:

* `desired_capacity` - (Required) The desired number of CVM instances that should be running in the group.
* `max_size` - (Required) The maximum size for the Auto Scaling group.
* `min_size` - (Required) The minimum size for the Auto Scaling group.
* `scaling_group_id` - (Required, ForceNew) ID of a scaling group.
* `schedule_action_name` - (Required) The name of this scaling action.
* `start_time` - (Required) The time for this action to start, in "YYYY-MM-DDThh:mm:ss+08:00" format (UTC+8).
* `end_time` - (Optional) The time for this action to end, in "YYYY-MM-DDThh:mm:ss+08:00" format (UTC+8).
* `recurrence` - (Optional) The time when recurring future actions will start. Start time is specified by the user following the Unix cron syntax format. And this argument should be set with end_time together.


