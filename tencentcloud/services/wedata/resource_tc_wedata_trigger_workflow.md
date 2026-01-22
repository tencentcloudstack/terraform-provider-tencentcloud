Provides a resource to create a wedata trigger workflow

Example Usage

```hcl
resource "tencentcloud_wedata_trigger_workflow" "workflow" {
  bundle_id          = null
  bundle_info        = null
  owner_uin          = 100044349576
  parent_folder_path = "/默认文件夹"
  project_id         = 3108707295180644352
  workflow_desc      = null
  workflow_name      = "tf-test1"
  general_task_params {
    type  = "SPARK_SQL"
    value = "a=b\nb=c\nc=d\nd=e"
  }
  trigger_workflow_scheduler_configurations {
    config_mode                     = "COMMON"
    crontab_expression              = "0 0 * * * ? *"
    cycle_type                      = "DAY_CYCLE"
    end_time                        = "2099-12-31 23:59:59"
    extra_info                      = null
    file_arrival_path               = null
    schedule_time_zone              = "UTC+8"
    scheduler_status                = "ACTIVE"
    start_time                      = "2026-01-09 00:00:00"
    trigger_minimum_interval_second = 0
    trigger_mode                    = "TIME_TRIGGER"
    trigger_wait_time_second        = 0
  }
  workflow_params {
    param_key   = "aaa"
    param_value = "bbb"
  }
  workflow_params {
    param_key   = "bbb"
    param_value = "ccc"
  }
}
```

Import

wedata trigger_workflow can be imported using the id, e.g.

```
terraform import tencentcloud_wedata_trigger_workflow.trigger_workflow project_id#workflow_id
```