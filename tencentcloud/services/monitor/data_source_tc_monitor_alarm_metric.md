Use this data source to query detailed information of monitor alarm_metric

Example Usage

```hcl
data "tencentcloud_monitor_alarm_metric" "alarm_metric" {
  module       = "monitor"
  monitor_type = "Monitoring"
  namespace    = "cvm_device"
}
```