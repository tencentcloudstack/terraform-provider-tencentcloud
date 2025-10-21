---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_async_event_status"
sidebar_current: "docs-tencentcloud-datasource-scf_async_event_status"
description: |-
  Use this data source to query detailed information of scf async_event_status
---

# tencentcloud_scf_async_event_status

Use this data source to query detailed information of scf async_event_status

## Example Usage

```hcl
data "tencentcloud_scf_async_event_status" "async_event_status" {
  invoke_request_id = "9de9405a-e33a-498d-bb59-e80b7bed1191"
}
```

## Argument Reference

The following arguments are supported:

* `invoke_request_id` - (Required, String) ID of the async execution request.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Async event status.
  * `invoke_request_id` - Async execution request ID.
  * `status_code` - Request status code.
  * `status` - Async event status. Values: `RUNNING` (running); `FINISHED` (invoked successfully); `ABORTED` (invocation ended); `FAILED` (invocation failed).


