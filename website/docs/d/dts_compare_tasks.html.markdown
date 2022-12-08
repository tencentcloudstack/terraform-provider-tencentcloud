---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_compare_tasks"
sidebar_current: "docs-tencentcloud-datasource-dts_compare_tasks"
description: |-
  Use this data source to query detailed information of dts compareTasks
---

# tencentcloud_dts_compare_tasks

Use this data source to query detailed information of dts compareTasks

## Example Usage

```hcl
data "tencentcloud_dts_compare_tasks" "compareTasks" {
  job_id = ""
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String) job id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - compare task list.
  * `check_process` - compare check info.
    * `message` - message.
    * `percent` - progress info.
    * `status` - status.
    * `step_all` - all step counts.
    * `step_now` - current step number.
    * `step` - step info.
  * `compare_process` - compare processing info.
    * `message` - message.
    * `percent` - progress info.
    * `status` - status.
    * `step_all` - all step counts.
    * `step_now` - current step number.
    * `step` - step info.
  * `compare_task_id` - compare task id.
  * `conclusion` - conclusion.
  * `config` - config.
    * `object_items` - object items.
    * `object_mode` - object mode.
  * `created_at` - create time.
  * `finished_at` - finished time.
  * `job_id` - job id.
  * `started_at` - start time.
  * `status` - compare task status, optional value is created/readyRun/running/success/stopping/failed/canceled.
  * `task_name` - compare task name.


