Use this data source to query detailed information of mysql user_task

Example Usage

```hcl
data "tencentcloud_mysql_user_task" "user_task" {
  instance_id = "cdb-fitq5t9h"
  async_request_id = "f2fe828c-773af816-0a08f542-94bb2a9c"
  task_types = [5]
  task_status = [2]
  start_time_begin = "2017-12-31 10:40:01"
  start_time_end = "2017-12-31 10:40:01"
}
```