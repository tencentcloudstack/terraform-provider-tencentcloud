Provides a resource to create a emr emr_auto_scale_strategy

Example Usage

```hcl
resource "tencentcloud_emr_auto_scale_strategy" "emr_auto_scale_strategy" {
  instance_id   = "emr-rzrochgp"
  strategy_type = 2
  time_auto_scale_strategy {
    strategy_name    = "tf-test1"
    interval_time    = 100
    scale_action     = 1
    scale_num        = 1
    strategy_status  = 1
    retry_valid_time = 60
    repeat_strategy {
      repeat_type = "DAY"
      day_repeat {
        execute_at_time_of_day = "16:30:00"
        step                   = 1
      }
      expire = "2026-02-20 23:59:59"
    }
    grace_down_flag = false
    tags {
      tag_key   = "createBy"
      tag_value = "terraform"
    }
    config_group_assigned = "{\"HDFS-2.8.5\":-1,\"YARN-2.8.5\":-1}"
    measure_method        = "INSTANCE"
    terminate_policy      = "DEFAULT"
    soft_deploy_info      = [1, 2]
    service_node_info     = [7]
    priority              = 1
  }
}
```

Import

emr emr_auto_scale_strategy can be imported using the id, e.g.

```
terraform import tencentcloud_emr_auto_scale_strategy.emr_auto_scale_strategy emr_auto_scale_strategy_id
```
