Provides a resource to create a wedata integration_realtime_task

Example Usage

```hcl
resource "tencentcloud_wedata_integration_realtime_task" "example" {
  project_id  = "1612982498218618880"
  task_name   = "tf_example"
  task_mode   = "1"
  description = "description."
  sync_type   = 1
  task_info {
    incharge    = "100028439226"
    executor_id = "20230313175748567418"
    config {
      name  = "concurrency"
      value = "1"
    }
    config {
      name  = "TaskManager"
      value = "1"
    }
    config {
      name  = "JobManager"
      value = "1"
    }
    config {
      name  = "TolerateDirtyData"
      value = "0"
    }
    config {
      name  = "CheckpointingInterval"
      value = "1"
    }
    config {
      name  = "CheckpointingIntervalUnit"
      value = "min"
    }
    config {
      name  = "RestartStrategyFixedDelayAttempts"
      value = "-1"
    }
    config {
      name  = "ResourceAllocationType"
      value = "0"
    }
    config {
      name  = "TaskAlarmRegularList"
      value = "35"
    }
  }
}
```

Import

wedata integration_realtime_task can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_integration_realtime_task.example 1776563389209296896#h9d39630a-ae45-4460-90b2-0b093cbfef5d
```