Use this data source to query detailed information of redis instance_task_list

Example Usage

```hcl
data "tencentcloud_redis_instance_task_list" "instance_task_list" {
  instance_id = "crs-c1nl9rpv"
  instance_name = ""
  project_ids = [""]
  task_types = [""]
  begin_time = "2021-12-30 00:00:00"
  end_time = "2021-12-30 00:00:00"
  task_status = [""]
  result = [""]
  operate_uin = [""]
}
```