Use this data source to query detailed information of monitor statistic_data

Example Usage

```hcl
data "tencentcloud_monitor_statistic_data" "statistic_data" {
  module       = "monitor"
  namespace    = "QCE/TKE2"
  metric_names = ["cpu_usage"]
  conditions {
    key      = "tke_cluster_instance_id"
    operator = "="
    value    = ["cls-mw2w40s7"]
  }
}
```
