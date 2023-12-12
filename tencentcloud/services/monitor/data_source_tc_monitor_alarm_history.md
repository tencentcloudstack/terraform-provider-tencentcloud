Use this data source to query detailed information of monitor alarm_history

Example Usage

```hcl
data "tencentcloud_monitor_alarm_history" "alarm_history" {
  module        = "monitor"
  order         = "DESC"
  start_time    = 1696608000
  end_time      = 1697212799
  monitor_types = ["MT_QCE"]
  project_ids   = [0]
  namespaces {
    monitor_type = "CpuUsage"
    namespace    = "cvm_device"
  }
  policy_name = "terraform_test"
  content     = "CPU利用率 > 3%"
  policy_ids  = ["policy-iejtp4ue"]
}
```