---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_shipper_tasks"
sidebar_current: "docs-tencentcloud-datasource-cls_shipper_tasks"
description: |-
  Use this data source to query detailed information of cls shipper_tasks
---

# tencentcloud_cls_shipper_tasks

Use this data source to query detailed information of cls shipper_tasks

## Example Usage

```hcl
data "tencentcloud_cls_shipper_tasks" "shipper_tasks" {
  shipper_id = "dbde3c9b-ea16-4032-bc2a-d8fa65567a8e"
  start_time = 160749910700
  end_time   = 160749910800
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, Int) end time(ms).
* `shipper_id` - (Required, String) shipper id.
* `start_time` - (Required, Int) start time(ms).
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tasks` - .
  * `end_time` - end time(ms).
  * `message` - detail info.
  * `range_end` - end time of current task (ms).
  * `range_start` - start time of current task (ms).
  * `shipper_id` - shipper id.
  * `start_time` - start time(ms).
  * `status` - status of current shipper task.
  * `task_id` - task id.
  * `topic_id` - topic id.


