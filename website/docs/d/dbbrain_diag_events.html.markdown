---
subcategory: "TencentDB for DBbrain(dbbrain)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dbbrain_diag_events"
sidebar_current: "docs-tencentcloud-datasource-dbbrain_diag_events"
description: |-
  Use this data source to query detailed information of dbbrain diag_events
---

# tencentcloud_dbbrain_diag_events

Use this data source to query detailed information of dbbrain diag_events

## Example Usage

```hcl
data "tencentcloud_dbbrain_diag_events" "diag_events" {
  instance_ids = ["%s"]
  start_time   = "%s"
  end_time     = "%s"
  severities   = [1, 4, 5]
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) end time.
* `start_time` - (Required, String) start time.
* `instance_ids` - (Optional, Set: [`String`]) instance id list.
* `result_output_file` - (Optional, String) Used to save results.
* `severities` - (Optional, Set: [`Int`]) severity list, optional value is 1-fatal, 2-severity, 3-warning, 4-tips, 5-health.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - diag event list.
  * `diag_item` - diag item.
  * `diag_type` - diag type.
  * `end_time` - end time.
  * `event_id` - event id.
  * `instance_id` - instance id.
  * `metric` - metric.
  * `outline` - outline.
  * `region` - region.
  * `severity` - severity.
  * `start_time` - start time.


