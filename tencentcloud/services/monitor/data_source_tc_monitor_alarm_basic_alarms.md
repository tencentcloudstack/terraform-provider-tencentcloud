Use this data source to query detailed information of monitor basic_alarms

Example Usage

```hcl
data "tencentcloud_monitor_alarm_basic_alarms" "alarms" {
  module             = "monitor"
  start_time         = 1696990903
  end_time           = 1697098903
  occur_time_order   = "DESC"
  project_ids        = [0]
  view_names         = ["cvm_device"]
  alarm_status       = [1]
  instance_group_ids = [5497073]
  metric_names       = ["cpu_usage"]
}
```