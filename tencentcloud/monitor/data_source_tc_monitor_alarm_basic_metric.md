Use this data source to query detailed information of monitor basic_metric

Example Usage

```hcl
data "tencentcloud_monitor_alarm_basic_metric" "alarm_metric" {
  namespace   = "qce/cvm"
  metric_name = "WanOuttraffic"
  dimensions  = ["uuid"]
}
```