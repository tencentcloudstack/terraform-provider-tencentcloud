---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_sql_script_runs"
sidebar_current: "docs-tencentcloud-datasource-wedata_sql_script_runs"
description: |-
  Use this data source to query detailed information of WeData sql script runs
---

# tencentcloud_wedata_sql_script_runs

Use this data source to query detailed information of WeData sql script runs

## Example Usage

```hcl
data "tencentcloud_wedata_sql_script_runs" "example" {
  project_id = "1460947878944567296"
  script_id  = "971c1520-836f-41be-b13f-7a6c637317c8"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `script_id` - (Required, String) Script ID.
* `end_time` - (Optional, String) End time.
* `execute_user_uin` - (Optional, String) Execute user UIN.
* `job_id` - (Optional, String) Job ID.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) Search keyword.
* `start_time` - (Optional, String) Start time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Data exploration tasks.
  * `create_time` - Task creation time.
  * `end_time` - End time.
  * `job_execution_list` - Subtask list.
    * `collecting_total_result` - Whether collecting full results: default false, true means collecting full results, used for frontend polling.
    * `context_script_content` - Context SQL content.
    * `create_time` - Create time.
    * `end_time` - End time.
    * `execute_stage_info` - Execution phase.
    * `job_execution_id` - Subquery task ID.
    * `job_execution_name` - Subquery name.
    * `job_id` - Data exploration task ID.
    * `log_file_path` - Log file path.
    * `result_effect_count` - Number of rows affected by the task execution result.
    * `result_file_path` - Result file path.
    * `result_preview_count` - Number of rows for previewing the task execution results.
    * `result_preview_file_path` - Preview result file path.
    * `result_total_count` - Total number of rows in the task execution result.
    * `script_content_truncate` - Whether the script content is truncated.
    * `script_content` - Subquery SQL content.
    * `status` - Subquery status.
    * `time_cost` - Time consumed.
    * `update_time` - Update time.
  * `job_id` - Data exploration task ID.
  * `job_name` - Data exploration task name.
  * `job_type` - Job type.
  * `owner_uin` - Cloud owner account UIN.
  * `script_content_truncate` - Whether the script content is truncated.
  * `script_content` - Script content.
  * `script_id` - Script ID.
  * `status` - Task status.
  * `time_cost` - Time consumed.
  * `update_time` - Update time.
  * `user_uin` - Account UIN.


