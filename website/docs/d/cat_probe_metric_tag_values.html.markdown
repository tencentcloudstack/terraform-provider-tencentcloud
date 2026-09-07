---
subcategory: "Cloud Automated Testing(CAT)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cat_probe_metric_tag_values"
sidebar_current: "docs-tencentcloud-datasource-cat_probe_metric_tag_values"
description: |-
  Use this data source to query detailed information of cat probe_metric_tag_values
---

# tencentcloud_cat_probe_metric_tag_values

Use this data source to query detailed information of cat probe_metric_tag_values

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `analyze_task_type` - (Optional, String) Analysis of task type, supported types: `AnalyzeTaskType_Network`: network quality, `AnalyzeTaskType_Browse`: page performance, `AnalyzeTaskType_Transport`: port performance, `AnalyzeTaskType_UploadDownload`: file transport, `AnalyzeTaskType_MediaStream`: audiovisual experience.
* `filter` - (Optional, String) Filter conditions can be passed as a single filter or multiple parameters concatenated together, support regular matching.
* `filters` - (Optional, Set: [`String`]) Multiple condition filtering, supports combining multiple filtering conditions for query.
* `key` - (Optional, String) Dimension tag value, reference: `host`: task domain, `errorInfo`: status type, `area`: probe point area, `operator`: probe point operator, `taskId`: task ID.
* `result_output_file` - (Optional, String) Used to save results.
* `time_range` - (Optional, String) Time range.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `tag_value_set` - Tag value serialized string.


