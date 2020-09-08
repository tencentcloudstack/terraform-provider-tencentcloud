---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_logs"
sidebar_current: "docs-tencentcloud-datasource-scf_logs"
description: |-
  Use this data source to query SCF function logs.
---

# tencentcloud_scf_logs

Use this data source to query SCF function logs.

## Example Usage

```hcl
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  cos_bucket_name   = "scf-code-1234567890"
  cos_object_name   = "code.zip"
  cos_bucket_region = "ap-guangzhou"
}

data "tencentcloud_scf_logs" "foo" {
  function_name = tencentcloud_scf_function.foo.name
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required) Name of the SCF function to be queried.
* `end_time` - (Optional) The end time of the query, the format is `2017-05-16 20:00:00`, which can only be within one day from `start_time`.
* `invoke_request_id` - (Optional) Corresponding requestId when executing function.
* `limit` - (Optional) Number of logs, the default is `10000`, offset+limit cannot be greater than 10000.
* `namespace` - (Optional) Namespace of the SCF function to be queried.
* `offset` - (Optional) Log offset, default is `0`, offset+limit cannot be greater than 10000.
* `order_by` - (Optional) Sort the logs according to the following fields: `function_name`, `duration`, `mem_usage`, `start_time`, default `start_time`.
* `order` - (Optional) Order to sort the log, optional values `desc` and `asc`, default `desc`.
* `result_output_file` - (Optional) Used to save results.
* `ret_code` - (Optional) Use to filter log, optional value: `not0` only returns the error log. `is0` only returns the correct log. `TimeLimitExceeded` returns the log of the function call timeout. `ResourceLimitExceeded` returns the function call generation resource overrun log. `UserCodeException` returns logs of the user code error that occurred in the function call. Not passing the parameter means returning all logs.
* `start_time` - (Optional) The start time of the query, the format is `2017-05-16 20:00:00`, which can only be within one day from `end_time`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `logs` - An information list of logs. Each element contains the following attributes:
  * `bill_duration` - Function billing time, according to duration up to the last 100ms, unit is ms.
  * `duration` - Function execution time-consuming, unit is ms.
  * `function_name` - Name of the SCF function.
  * `invoke_finished` - Whether the function call ends, `1` means the execution ends, other values indicate the call exception.
  * `level` - Log level.
  * `log` - Log output during function execution.
  * `mem_usage` - The actual memory size consumed in the execution of the function, unit is Byte.
  * `request_id` - Execute the requestId corresponding to the function.
  * `ret_code` - Execution result of function, `0` means the execution is successful, other values indicate failure.
  * `ret_msg` - Return value after function execution is completed.
  * `source` - Log source.
  * `start_time` - Point in time at which the function begins execution.


