---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_tag_retention_execution_tasks"
sidebar_current: "docs-tencentcloud-datasource-tcr_tag_retention_execution_tasks"
description: |-
  Use this data source to query detailed information of tcr tag_retention_execution_tasks
---

# tencentcloud_tcr_tag_retention_execution_tasks

Use this data source to query detailed information of tcr tag_retention_execution_tasks

## Example Usage

```hcl
data "tencentcloud_tcr_tag_retention_execution_tasks" "tasks" {
  registry_id  = "%s"
  retention_id = "17"
  execution_id = "1"
}
```

## Argument Reference

The following arguments are supported:

* `execution_id` - (Required, Int) execution id.
* `registry_id` - (Required, String) instance id.
* `retention_id` - (Required, Int) retention id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `retention_task_list` - list of version retention tasks.
  * `end_time` - task end time.
  * `execution_id` - the rule execution id.
  * `repository` - repository name.
  * `retained` - Total number of retained tags.
  * `start_time` - task start time.
  * `status` - the execution status of the task: Failed, Succeed, Stopped, InProgress.
  * `task_id` - task id.
  * `total` - Total number of tags.


