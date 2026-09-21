---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_inference_service_deployment_logs"
sidebar_current: "docs-tencentcloud-datasource-teo_inference_service_deployment_logs"
description: |-
  Use this data source to query detailed information of teo inference service deployment logs
---

# tencentcloud_teo_inference_service_deployment_logs

Use this data source to query detailed information of teo inference service deployment logs

## Example Usage

```hcl
data "tencentcloud_teo_inference_service_deployment_logs" "example" {
  zone_id    = "zone-2qtuhspy7cr6"
  service_id = "sid-2quhspyeq8r6"
  record_id  = "rec-2quhspyeq8r6"
}
```

## Argument Reference

The following arguments are supported:

* `record_id` - (Required, String) Deployment record ID.
* `service_id` - (Required, String) Inference service ID.
* `zone_id` - (Required, String) Site ID.
* `end_time` - (Optional, String) End time of the logs to be retrieved. The default time range (EndTime - StartTime) is the last 7 days.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) Sort field. Valid value: `timestamp` (log generation time). Default value: `timestamp`.
* `sort_order` - (Optional, String) Sort order. Valid values: `asc` (ascending), `desc` (descending). Default value: `desc`.
* `start_time` - (Optional, String) Start time of the logs to be retrieved.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `deployment_log_info_set` - Deployment log list.
  * `log_message` - Log message content.
  * `timestamp` - Log generation time.


