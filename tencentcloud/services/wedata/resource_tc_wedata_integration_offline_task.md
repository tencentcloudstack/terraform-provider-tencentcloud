Provides a resource to create a wedata integration_offline_task

Example Usage

```hcl
resource "tencentcloud_wedata_integration_offline_task" "example" {
  project_id  = "1612982498218618880"
  cycle_step  = 1
  delay_time  = 0
  end_time    = "2099-12-31 00:00:00"
  notes       = "terraform example demo."
  start_time  = "2023-12-31 00:00:00"
  task_name   = "tf_example"
  task_action = "2"
  task_mode   = "1"

  task_info {
    executor_id = "20230313175748567418"
    config {
      name  = "Args"
      value = "args"
    }
    config {
      name  = "dirtyDataThreshold"
      value = "0"
    }
    config {
      name  = "concurrency"
      value = "1"
    }
    config {
      name  = "syncRateLimitUnit"
      value = "0"
    }
    ext_config {
      name  = "TaskAlarmRegularList"
      value = "73"
    }
    incharge = "demo"
    offline_task_add_entity {
      cycle_type         = 3
      crontab_expression = "0 0 1 * * ?"
      retry_wait         = 5
      retriable          = 1
      try_limit          = 5
      self_depend        = 1
    }
  }
}
```

Import

wedata integration_offline_task can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_integration_offline_task.example 1612982498218618880#20231102200955095
```