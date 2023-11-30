Use this data source to query detailed information of cat metric_data
Example Usage
```hcl
data "tencentcloud_cat_metric_data" "metric_data" {
  analyze_task_type = "AnalyzeTaskType_Network"
  metric_type = "gauge"
  field = "avg(\"ping_time\")"
  filters = [
    "\"host\" = 'www.qq.com'",
    "time >= now()-1h",
  ]
}
```