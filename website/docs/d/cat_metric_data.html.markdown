---
subcategory: "Cloud Automated Testing(CAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cat_metric_data"
sidebar_current: "docs-tencentcloud-datasource-cat_metric_data"
description: |-
  Use this data source to query detailed information of cat metric_data
---

# tencentcloud_cat_metric_data

Use this data source to query detailed information of cat metric_data

## Example Usage

```hcl
data "tencentcloud_cat_metric_data" "metric_data" {
  analyze_task_type = "AnalyzeTaskType_Network"
  metric_type       = "gauge"
  field             = "avg(\"ping_time\")"
  filters = [
    "\"host\" = 'www.qq.com'",
    "time >= now()-1h",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `analyze_task_type` - (Required, String) Analysis of task type, supported types: `AnalyzeTaskType_Network`: network quality, `AnalyzeTaskType_Browse`: page performance, `AnalyzeTaskType_Transport`: port performance, `AnalyzeTaskType_UploadDownload`: file transport, `AnalyzeTaskType_MediaStream`: audiovisual experience.
* `field` - (Required, String) Detailed fields of metrics, specified metrics can be passed or aggregate metrics, such as avg(ping_time) means entire delay.
* `filters` - (Required, Set: [`String`]) Multiple condition filtering, supports combining multiple filtering conditions for query.
* `metric_type` - (Required, String) Metric type, metrics queries are passed with gauge by default.
* `filter` - (Optional, String) Filter conditions can be passed as a single filter or multiple parameters concatenated together.
* `group_by` - (Optional, String) Aggregation time, such as 1m, 1d, 30d, and so on.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `metric_set` - Return JSON string.


