Use this data source to query TEO inference service monitor data, such as CPU usage, GPU usage, instance count, and memory usage metrics.

Example Usage

```hcl
data "tencentcloud_teo_inference_service_monitor_data" "basic" {
  zone_id      = "zone-xxxxx"
  service_ids  = ["service-xxxxx"]
  metric_names = ["cpu_usage_average"]
  start_time   = "2024-01-01T00:00:00Z"
  end_time     = "2024-01-01T01:00:00Z"
}

output "records" {
  value = data.tencentcloud_teo_inference_service_monitor_data.basic.inference_service_monitor_records
}
```

```hcl
data "tencentcloud_teo_inference_service_monitor_data" "with_interval" {
  zone_id      = "zone-xxxxx"
  service_ids  = ["service-xxxxx"]
  metric_names = ["cpu_usage_average", "gpu_usage_max"]
  start_time   = "2024-01-01T00:00:00Z"
  end_time     = "2024-01-02T00:00:00Z"
  interval      = "hour"
}

output "records" {
  value = data.tencentcloud_teo_inference_service_monitor_data.with_interval.inference_service_monitor_records
}
```

```hcl
data "tencentcloud_teo_inference_service_monitor_data" "export" {
  zone_id            = "zone-xxxxx"
  service_ids         = ["service-xxxxx"]
  metric_names        = ["cpu_usage_average"]
  start_time          = "2024-01-01T00:00:00Z"
  end_time            = "2024-01-01T01:00:00Z"
  result_output_file = "./monitor_data.json"
}
```