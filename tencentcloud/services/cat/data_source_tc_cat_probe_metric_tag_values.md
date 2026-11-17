Use this data source to query detailed information of cat probe_metric_tag_values
Example Usage
```hcl
data "tencentcloud_cat_probe_metric_tag_values" "example" {
  analyze_task_type = "AnalyzeTaskType_Network"
  key               = "host"
  filter            = "www.qq.com"
  time_range        = "1h"
}

output "tag_values" {
  value = data.tencentcloud_cat_probe_metric_tag_values.example.tag_value_set
}
```
