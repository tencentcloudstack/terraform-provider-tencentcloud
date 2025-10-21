---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_async_event_management"
sidebar_current: "docs-tencentcloud-datasource-scf_async_event_management"
description: |-
  Use this data source to query detailed information of scf async_event_management
---

# tencentcloud_scf_async_event_management

Use this data source to query detailed information of scf async_event_management

## Example Usage

```hcl
data "tencentcloud_scf_async_event_management" "async_event_management" {
  function_name = "keep-1676351130"
  namespace     = "default"
  qualifier     = "$LATEST"
  order         = "ASC"
  orderby       = "StartTime"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String) Function name.
* `invoke_request_id` - (Optional, String) Filter (event invocation request ID).
* `invoke_type` - (Optional, Set: [`String`]) Filter (invocation type list), Values: CMQ, CKAFKA_TRIGGER, APIGW, COS, TRIGGER_TIMER, MPS_TRIGGER, CLS_TRIGGER, OTHERS.
* `namespace` - (Optional, String) Function namespace.
* `order` - (Optional, String) Valid values: ASC, DESC. Default value: DESC.
* `orderby` - (Optional, String) Valid values: StartTime, EndTime. Default value: StartTime.
* `qualifier` - (Optional, String) Filter (function version).
* `result_output_file` - (Optional, String) Used to save results.
* `status` - (Optional, Set: [`String`]) Filter (event status list), Values: RUNNING, FINISHED, ABORTED, FAILED.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `event_list` - Async event list.
  * `end_time` - Invocation end time in the format of %Y-%m-%d %H:%M:%S.%f.
  * `invoke_request_id` - Invocation request ID.
  * `invoke_type` - Invocation type.
  * `qualifier` - Function version.
  * `start_time` - Invocation start time in the format of %Y-%m-%d %H:%M:%S.%f.
  * `status` - Event status. Values: `RUNNING`; `FINISHED` (invoked successfully); `ABORTED` (invocation ended); `FAILED` (invocation failed).


