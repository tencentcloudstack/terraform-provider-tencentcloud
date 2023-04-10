---
subcategory: "TencentDB for Redis(crs)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_redis_instance_task_list"
sidebar_current: "docs-tencentcloud-datasource-redis_instance_task_list"
description: |-
  Use this data source to query detailed information of redis instance_task_list
---

# tencentcloud_redis_instance_task_list

Use this data source to query detailed information of redis instance_task_list

## Example Usage

```hcl
data "tencentcloud_redis_instance_task_list" "instance_task_list" {
  instance_id   = "crs-c1nl9rpv"
  instance_name = ""
  project_ids   = [""]
  task_types    = [""]
  begin_time    = "2021-12-30 00:00:00"
  end_time      = "2021-12-30 00:00:00"
  task_status   = [""]
  result        = [""]
  operate_uin   = [""]
}
```

## Argument Reference

The following arguments are supported:

* `begin_time` - (Optional, String) Start time.
* `end_time` - (Optional, String) Termination time.
* `instance_id` - (Optional, String) The ID of instance.
* `instance_name` - (Optional, String) Instance name.
* `operate_uin` - (Optional, Set: [`String`]) Operator Uin.
* `project_ids` - (Optional, Set: [`Int`]) Project Id.
* `result_output_file` - (Optional, String) Used to save results.
* `result` - (Optional, Set: [`Int`]) Task status.
* `task_status` - (Optional, Set: [`Int`]) Task status.
* `task_types` - (Optional, Set: [`String`]) Task type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tasks` - Task details.
  * `end_time` - The end time.
  * `instance_id` - The ID of instance.
  * `instance_name` - The name of instance.
  * `progress` - Task progress.
  * `project_id` - The project ID.
  * `result` - Task status.
  * `start_time` - Start time.
  * `task_id` - Task ID.
  * `task_type` - Task type.


