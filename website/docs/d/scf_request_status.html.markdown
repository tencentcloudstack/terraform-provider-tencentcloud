---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_request_status"
sidebar_current: "docs-tencentcloud-datasource-scf_request_status"
description: |-
  Use this data source to query detailed information of scf request_status
---

# tencentcloud_scf_request_status

Use this data source to query detailed information of scf request_status

## Example Usage

```hcl
data "tencentcloud_scf_request_status" "request_status" {
  function_name       = "keep-1676351130"
  function_request_id = "9de9405a-e33a-498d-bb59-e80b7bed1191"
  namespace           = "default"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String) Function name.
* `function_request_id` - (Required, String) ID of the request to be queried.
* `end_time` - (Optional, String) End time of the query. such as `2017-05-16 20:59:59`. If `StartTime` is not specified, `EndTime` defaults to the current time. If `StartTime` is specified, `EndTime` is required, and it need to be later than the `StartTime`.
* `namespace` - (Optional, String) Function namespace.
* `result_output_file` - (Optional, String) Used to save results.
* `start_time` - (Optional, String) Start time of the query, for example `2017-05-16 20:00:00`. If it's left empty, it defaults to 15 minutes before the current time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Details of the function running statusNote: this field may return `null`, indicating that no valid values can be obtained.
  * `duration` - Time consumed for the request in ms.
  * `function_name` - Function name.
  * `mem_usage` - Time consumed by the request in MB.
  * `request_id` - Request ID.
  * `ret_code` - Result of the request. `0`: succeeded, `1`: running, `-1`: exception.
  * `ret_msg` - Return value after the function is executed.
  * `retry_num` - Retry Attempts.
  * `start_time` - Request start time.


