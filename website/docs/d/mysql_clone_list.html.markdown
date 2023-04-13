---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_clone_list"
sidebar_current: "docs-tencentcloud-datasource-mysql_clone_list"
description: |-
  Use this data source to query detailed information of mysql clone_list
---

# tencentcloud_mysql_clone_list

Use this data source to query detailed information of mysql clone_list

## Example Usage

```hcl
data "tencentcloud_mysql_clone_list" "clone_list" {
  instance_id = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Query the list of cloning tasks for the specified source instance.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Clone task list.
  * `clone_job_id` - Id of the task list corresponding to the clone task.
  * `dst_instance_id` - The newly spawned instance Id of the clone task.
  * `end_time` - Task end time.
  * `new_region_id` - Id of the region where the cloned instance is located.
  * `rollback_strategy` - The strategy used by the clone instance includes the following types: timepoint: specify the point-in-time rollback, backupset: specify the backup file rollback.
  * `rollback_target_time` - The time point when the clone instance is rolled back.
  * `src_instance_id` - The source instance Id of the clone task.
  * `src_region_id` - Id of the region where the source instance is located.
  * `start_time` - Task start time.
  * `task_status` - Task status, including the following status: initial, running, wait_complete, success, failed.


