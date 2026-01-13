Provides a resource to create a wedata quality rule group

Example Usage

```hcl
resource "tencentcloud_wedata_quality_rule_group" "quality_rule_group" {
  project_id = "1612982498218618880"
  
  rule_group_exec_strategy_bo_list {
    monitor_type        = 3
    executor_group_id   = "20220509_114220_111"
    rule_group_name     = "tf-test-quality-rule-group"
    database_name       = "tf_test_db"
    datasource_id       = "20220509_114220_111"
    table_name          = "tf_test_table"
    exec_queue          = "default"
    executor_group_name = "default-group"
    
    start_time   = "2024-01-01 00:00:00"
    end_time     = "2024-12-31 23:59:59"
    cycle_type   = "D"
    delay_time   = 0
    cycle_step   = 1
    
    description          = "Test quality rule group"
    schedule_time_zone   = "UTC+8"
    exec_engine_type     = "SPARK"
    
    group_config {
      analysis_type = "SNAPSHOT"
    }
  }
}
```
