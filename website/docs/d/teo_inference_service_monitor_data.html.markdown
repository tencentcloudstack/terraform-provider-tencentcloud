---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_service_monitor_data"
sidebar_current: "docs-tencentcloud-datasource-teo_inference_service_monitor_data"
description: |-
  Use this data source to query TEO inference service monitor data, such as CPU usage, GPU usage, instance count, and memory usage metrics.
---

# tencentcloud_teo_inference_service_monitor_data

Use this data source to query TEO inference service monitor data, such as CPU usage, GPU usage, instance count, and memory usage metrics.

## Example Usage

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
  interval     = "hour"
}

output "records" {
  value = data.tencentcloud_teo_inference_service_monitor_data.with_interval.inference_service_monitor_records
}
```



```hcl
data "tencentcloud_teo_inference_service_monitor_data" "export" {
  zone_id            = "zone-xxxxx"
  service_ids        = ["service-xxxxx"]
  metric_names       = ["cpu_usage_average"]
  start_time         = "2024-01-01T00:00:00Z"
  end_time           = "2024-01-01T01:00:00Z"
  result_output_file = "./monitor_data.json"
}
```

## Argument Reference

The following arguments are supported:

* `end_time` - (Required, String) End time. The query time range (`EndTime` - `StartTime`) must be less than or equal to 30 days.
* `metric_names` - (Required, List: [`String`]) Metric list, up to 10 metrics. Valid values: `cpu_usage_average`, `cpu_usage_max`, `gpu_usage_average`, `gpu_usage_max`, `instance_num_average`, `instance_num_max`, `gpu_memory_usage_max`, `memory_usage_average`, `memory_usage_max`.
* `service_ids` - (Required, List: [`String`]) Inference service ID list, up to 10 inference service IDs.
* `start_time` - (Required, String) Start time.
* `zone_id` - (Required, String) Site ID.
* `interval` - (Optional, String) Query time granularity. Valid values: `min`, `5min`, `hour`, `day`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `inference_service_monitor_records` - Inference service monitor record list.
  * `inference_service_monitor_items` - Detailed inference service monitor data.
    * `timestamp` - Timestamp of the monitor data point.
    * `value` - Numeric value of the monitor data point.
  * `metric_name` - Metric name.
  * `service_id` - Inference service ID.


