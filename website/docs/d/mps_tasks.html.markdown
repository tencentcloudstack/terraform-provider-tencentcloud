---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_tasks"
sidebar_current: "docs-tencentcloud-datasource-mps_tasks"
description: |-
  Use this data source to query detailed information of mps tasks
---

# tencentcloud_mps_tasks

Use this data source to query detailed information of mps tasks

## Example Usage

```hcl
data "tencentcloud_mps_tasks" "tasks" {
  status = "FINISH"
  limit  = 20
}
```

## Argument Reference

The following arguments are supported:

* `status` - (Required, String) Filter condition: task status, optional values: WAITING, PROCESSING, FINISH.
* `limit` - (Optional, Int) Return the number of records, default value: 10, maximum value: 100.
* `result_output_file` - (Optional, String) Used to save results.
* `scroll_token` - (Optional, String) Page turning flag, used when pulling in batches: when a single request cannot pull all the data, the interface will return a ScrollToken, and the next request will carry this Token, and it will be obtained from the next record.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `task_set` - Task list.
  * `begin_process_time` - Begin process time, in ISO date format. Refer to https://cloud.tencent.com/document/product/862/37710#52. If the task has not started yet, this field is: 0000-00-00T00:00:00Z.
  * `create_time` - Creation time, in ISO date format. Refer to https://cloud.tencent.com/document/product/862/37710#52.
  * `finish_time` - Task finish time, in ISO date format. Refer to https://cloud.tencent.com/document/product/862/37710#52. If the task has not been completed, this field is: 0000-00-00T00:00:00Z.
  * `sub_task_types` - Sub task types.
  * `task_id` - Task ID.
  * `task_type` - Task type, including:WorkflowTask, EditMediaTask, LiveProcessTask.


