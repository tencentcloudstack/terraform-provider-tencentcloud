---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_user_task"
sidebar_current: "docs-tencentcloud-datasource-mysql_user_task"
description: |-
  Use this data source to query detailed information of mysql user_task
---

# tencentcloud_mysql_user_task

Use this data source to query detailed information of mysql user_task

## Example Usage

```hcl
data "tencentcloud_mysql_user_task" "user_task" {
  instance_id      = "cdb-fitq5t9h"
  async_request_id = "f2fe828c-773af816-0a08f542-94bb2a9c"
  task_types       = [5]
  task_status      = [2]
  start_time_begin = "2017-12-31 10:40:01"
  start_time_end   = "2017-12-31 10:40:01"
}
```

## Argument Reference

The following arguments are supported:

* `async_request_id` - (Optional, String) Asynchronous task request ID, the AsyncRequestId returned by executing cloud database-related operations.
* `instance_id` - (Optional, String) Instance ID, the format is: cdb-c1nl9rpv, which is the same as the instance ID displayed on the cloud database console page, and you can use the [query instance list] (https://cloud.tencent.com/document/api/236/15872) interface Gets the value of the field InstanceId in the output parameter.
* `result_output_file` - (Optional, String) Used to save results.
* `start_time_begin` - (Optional, String) The start time of the first task, used for range query, the time format is as follows: 2017-12-31 10:40:01.
* `start_time_end` - (Optional, String) The start time of the last task, used for range query, the time format is as follows: 2017-12-31 10:40:01.
* `task_status` - (Optional, Set: [`String`]) Task status. If no value is passed, all task statuses will be queried. Supported values include: `UNDEFINED` - undefined; `INITIAL` - initialization; `RUNNING` - running; `SUCCEED` - the execution was successful; `FAILED` - execution failed; `KILLED` - terminated; `REMOVED` - removed; `PAUSED` - Paused.
* `task_types` - (Optional, Set: [`String`]) Task type. If no value is passed, all task types will be queried. Supported values include: `ROLLBACK` - database rollback; `SQL OPERATION` - SQL operation; `IMPORT DATA` - data import; `MODIFY PARAM` - parameter setting; `INITIAL` - initialize the cloud database instance; `REBOOT` - restarts the cloud database instance; `OPEN GTID` - open the cloud database instance GTID; `UPGRADE RO` - read-only instance upgrade; `BATCH ROLLBACK` - database batch rollback; `UPGRADE MASTER` - master upgrade; `DROP TABLES` - delete cloud database tables; `SWITCH DR TO MASTER` - The disaster recovery instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - The returned instance task information.
  * `async_request_id` - The request ID of the asynchronous task.
  * `code` - error code.
  * `end_time` - Instance task end time.
  * `instance_ids` - The instance ID associated with the task. Note: This field may return null, indicating that no valid value can be obtained.
  * `job_id` - Instance task ID.
  * `message` - error message.
  * `progress` - Instance task progress.
  * `start_time` - Instance task start time.
  * `task_status` - Instance task status, possible values include:UNDEFINED - undefined;INITIAL - initialization;RUNNING - running;SUCCEED - the execution was successful;FAILED - execution failed;KILLED - terminated;REMOVED - removed;PAUSED - Paused.WAITING - waiting (cancellable).
  * `task_type` - Instance task type, possible values include:ROLLBACK - database rollback;SQL OPERATION - SQL operation;IMPORT DATA - data import;MODIFY PARAM - parameter setting;INITIAL - initialize the cloud database instance;REBOOT - restarts the cloud database instance;OPEN GTID - open the cloud database instance GTID;UPGRADE RO - read-only instance upgrade;BATCH ROLLBACK - database batch rollback;UPGRADE MASTER - master upgrade;DROP TABLES - delete cloud database tables;SWITCH DR TO MASTER - The disaster recovery instance.


