---
subcategory: "Elasticsearch Service(ES)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_elasticsearch_instance_operations"
sidebar_current: "docs-tencentcloud-datasource-elasticsearch_instance_operations"
description: |-
  Use this data source to query detailed information of elasticsearch instance operations
---

# tencentcloud_elasticsearch_instance_operations

Use this data source to query detailed information of elasticsearch instance operations

## Example Usage

```hcl
data "tencentcloud_elasticsearch_instance_operations" "instance_operations" {
  instance_id = "es-xxxxxx"
  start_time  = "2018-01-01 00:00:00"
  end_time    = "2023-10-31 10:12:45"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time, e.g. 2019-03-30 20:18:03.
* `instance_id` - (Required, String) Instance id.
* `start_time` - (Required, String) Start time, e.g. 2019-03-07 16:30:39.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `operations` - Operation records.
  * `detail` - Operation details.
    * `new_info` - Configuration information after instance update.
      * `key` - Key.
      * `value` - Value.
    * `old_info` - Instance original configuration information.
      * `key` - Key.
      * `value` - Value.
  * `id` - Id.
  * `progress` - Operation progress.
  * `result` - Operation result.
  * `start_time` - Start time.
  * `sub_account_uin` - Operator uin.
  * `tasks` - Task information.
    * `elapsed_time` - Elapsed time.
    * `finish_time` - Task completion time.
    * `name` - Task name.
    * `process_info` - Progress info.
      * `completed` - Completed quantity.
      * `remain` - Remaining quantity.
      * `task_type` - Task type. 60: restart task 70: fragment migration task 80: node modification task.
      * `total` - Total quantity.
    * `progress` - Task progress.
    * `sub_tasks` - Subtask.
      * `err_msg` - Subtask error message.
      * `failed_indices` - The index name of the failed upgrade check.
      * `finish_time` - Subtask end time.
      * `level` - Subtask level, 1: warning; 2: failed.
      * `name` - Subtask name.
      * `result` - Subtask result.
      * `status` - Subtask status, 1: success; 0: processing; -1: failure.
      * `type` - Subtask type.
  * `type` - Type.


