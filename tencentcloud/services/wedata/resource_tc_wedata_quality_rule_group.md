Provides a resource to create a wedata quality rule group

Example Usage

```hcl
resource "tencentcloud_wedata_quality_rule_group" "group" {
  project_id = 3016337760439783424
  rule_group_exec_strategy_bo_list {
    cycle_step          = 0
    database_name       = "default"
    datasource_id       = 65253
    delay_time          = 0
    description         = "tf测试"
    exec_engine_type    = "HIVE"
    exec_plan           = null
    exec_queue          = "default"
    executor_group_id   = 20250807142245848024
    executor_group_name = "重庆调度资源组-2a8lsema"
    monitor_type        = 2
    rule_group_name     = "tf_test_2"
    rule_id             = 0
    table_name          = "big_table_500"
    trigger_types       = ["CYCLE", "MAKE_UP"]
    tasks {
      cycle_type          = 0
      in_charge_id_list   = []
      in_charge_name_list = ["hannahlliao"]
      schedule_time_zone  = null
      task_id             = 20251118145318149
      task_name           = "hannah_test111"
      task_type           = 2
      workflow_id         = "DATA_INTEGRATION_2025-11-01_1"
    }
  }
}
```
