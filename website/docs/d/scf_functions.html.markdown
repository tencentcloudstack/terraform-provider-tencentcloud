---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_functions"
sidebar_current: "docs-tencentcloud-datasource-scf_functions"
description: |-
  Use this data source to query SCF functions.
---

# tencentcloud_scf_functions

Use this data source to query SCF functions.

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

data "tencentcloud_scf_functions" "foo" {
  name = tencentcloud_scf_function.foo.name
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) Description of the SCF function to be queried.
* `name` - (Optional) Name of the SCF function to be queried.
* `namespace` - (Optional) Namespace of the SCF function to be queried.
* `result_output_file` - (Optional) Used to save results.
* `tags` - (Optional) Tags of the SCF function to be queried, can use up to 10 tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `functions` - An information list of functions. Each element contains the following attributes:
  * `cls_logset_id` - CLS logset ID of the SCF function.
  * `cls_topic_id` - CLS topic ID of the SCF function.
  * `code_error` - Code error of the SCF function.
  * `code_result` - Code result of the SCF function.
  * `code_size` - Code size of the SCF function.
  * `create_time` - Create time of the SCF function.
  * `description` - Description of the SCF function.
  * `eip_fixed` - Whether EIP is a fixed IP.
  * `eips` - EIP list of the SCF function.
  * `environment` - Environment variable of the SCF function.
  * `err_no` - Errno of the SCF function.
  * `handler` - Handler of the SCF function.
  * `host` - Host of the SCF function.
  * `install_dependency` - Whether to automatically install dependencies.
  * `l5_enable` - Whether to enable L5.
  * `mem_size` - Memory size of the SCF function runtime, unit is M.
  * `modify_time` - Modify time of the SCF function.
  * `name` - Name of the SCF function.
  * `namespace` - Namespace of the SCF function.
  * `role` - CAM role of the SCF function.
  * `runtime` - Runtime of the SCF function.
  * `status_desc` - Status description of the SCF function.
  * `status` - Status of the SCF function.
  * `subnet_id` - Subnet ID of the SCF function.
  * `tags` - Tags of the SCF function.
  * `timeout` - Timeout of the SCF function maximum execution time, unit is second.
  * `trigger_info` - Trigger details list the SCF function. Each element contains the following attributes:
    * `create_time` - Create time of the SCF function trigger.
    * `custom_argument` - user-defined parameter of the SCF function trigger.
    * `enable` - Whether to enable SCF function trigger.
    * `modify_time` - Modify time of the SCF function trigger.
    * `name` - Name of the SCF function trigger.
    * `trigger_desc` - TriggerDesc of the SCF function trigger.
    * `type` - Type of the SCF function trigger.
  * `vip` - Vip of the SCF function.
  * `vpc_id` - VPC ID of the SCF function.


