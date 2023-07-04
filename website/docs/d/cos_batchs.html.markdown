---
subcategory: "Cloud Object Storage(COS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cos_batchs"
sidebar_current: "docs-tencentcloud-datasource-cos_batchs"
description: |-
  Use this data source to query the COS batch.
---

# tencentcloud_cos_batchs

Use this data source to query the COS batch.

## Example Usage

```hcl
data "tencentcloud_cos_batchs" "cos_batchs" {
  uin   = "xxxxxx"
  appid = "xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `appid` - (Required, Int) Appid.
* `uin` - (Required, String) Uin.
* `job_statuses` - (Optional, String) The task status information you need to query. If you do not specify a task status, COS returns the status of all tasks that have been executed, including those that are in progress. If you specify a task status, COS returns the task in the specified state. Optional task states include: Active, Cancelled, Cancelling, Complete, Completing, Failed, Failing, New, Paused, Pausing, Preparing, Ready, Suspended.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `jobs` - Multiple batch processing task information.
  * `creation_time` - Job creation time.
  * `description` - Mission description. The length is limited to 0-256 bytes.
  * `job_id` - Job ID. The length is limited to 1-64 bytes.
  * `operation` - Actions performed on objects in a batch processing job. For example, COSPutObjectCopy.
  * `priority` - Mission priority. Tasks with higher values will be given priority. The priority size is limited to 0-2147483647.
  * `progress_summary` - Summary of the status of task implementation. Describe the total number of operations performed in this task, the number of successful operations, and the number of failed operations.
    * `number_of_tasks_failed` - The current failed Operand.
    * `number_of_tasks_succeeded` - The current successful Operand.
    * `total_number_of_tasks` - Total Operand.
  * `status` - Task execution status. Legal parameter values include Active, Cancelled, Cancelling, Complete, Completing, Failed, Failing, New, Paused, Pausing, Preparing, Ready, Suspended.
  * `termination_date` - The end time of the batch processing job.


