---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_function"
sidebar_current: "docs-tencentcloud-resource-scf_function"
description: |-
  Provide a resource to create a SCF function.
---

# tencentcloud_scf_function

Provide a resource to create a SCF function.

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
```

## Argument Reference

The following arguments are supported:

* `handler` - (Required) Handler of the SCF function. The format of name is `<filename>.<method_name>`, and it supports 26 English letters, numbers, connectors, and underscores, it should start with a letter. The last character cannot be `-` or `_`. Available length is 2-60.
* `name` - (Required, ForceNew) Name of the SCF function. Name supports 26 English letters, numbers, connectors, and underscores, it should start with a letter. The last character cannot be `-` or `_`. Available length is 2-60.
* `runtime` - (Required) Runtime of the SCF function, only supports `Python2.7`, `Python3.6`, `Nodejs6.10`, `Nodejs8.9`, `Nodejs10.15`, `PHP5`, `PHP7`, `Golang1`, and `Java8`.
* `cls_logset_id` - (Optional) cls logset id of the SCF function.
* `cls_topic_id` - (Optional) cls topic id of the SCF function.
* `cos_bucket_name` - (Optional) Cos bucket name of the SCF function, such as `cos-1234567890`, conflict with `zip_file`.
* `cos_bucket_region` - (Optional) Cos bucket region of the SCF function, conflict with `zip_file`.
* `cos_object_name` - (Optional) Cos object name of the SCF function, should have suffix `.zip` or `.jar`, conflict with `zip_file`.
* `description` - (Optional) Description of the SCF function. Description supports English letters, numbers, spaces, commas, newlines, periods and Chinese, the maximum length is 1000.
* `environment` - (Optional) Environment of the SCF function.
* `l5_enable` - (Optional) Enable L5 for SCF function, default is `false`.
* `mem_size` - (Optional) Memory size of the SCF function, unit is MB. The default is `128`MB. The range is 128M-1536M, and the ladder is 128M.
* `namespace` - (Optional, ForceNew) Namespace of the SCF function, default is `default`.
* `role` - (Optional) Role of the SCF function.
* `subnet_id` - (Optional) Subnet id of the SCF function.
* `tags` - (Optional) Tags of the SCF function.
* `timeout` - (Optional) Timeout of the SCF function, unit is second. Default `3`. Available value is 1-300.
* `triggers` - (Optional) Trigger list of the SCF function, note that if you modify the trigger list, all existing triggers will be deleted, and then create triggers in the new list. Each element contains the following attributes:
* `vpc_id` - (Optional) VPC id of the SCF function.
* `zip_file` - (Optional) Zip file of the SCF function, conflict with `cos_bucket_name`, `cos_object_name`, `cos_bucket_region`.

The `triggers` object supports the following:

* `name` - (Required) Name of the SCF function trigger, if `type` is `ckafka`, the format of name must be `<ckafkaInstanceId>-<topicId>`; if `type` is `cos`, the name is cos bucket id, other In any case, it can be combined arbitrarily. It can only contain English letters, numbers, connectors and underscores. The maximum length is 100.
* `trigger_desc` - (Required) TriggerDesc of the SCF function trigger, parameter format of `timer` is linux cron expression; parameter of `cos` type is json string `{"event":"cos:ObjectCreated:*","filter":{"Prefix":"","Suffix":""}}`, where `event` is the cos event trigger, `Prefix` is the corresponding file prefix filter condition, `Suffix` is the suffix filter condition, if not need filter condition can not pass; `cmq` type does not pass this parameter; `ckafka` type parameter format is json string `{"maxMsgNum":"1","offset":"latest"}`; `apigw` type parameter format is json string `{"api":{"authRequired":"FALSE","requestConfig":{"method":"ANY"},"isIntegratedResponse":"FALSE"},"service":{"serviceId":"service-dqzh68sg"},"release":{"environmentName":"test"}}`.
* `type` - (Required) Type of the SCF function trigger, support `cos`, `cmq`, `timer`, `ckafka`, `apigw`.
* `cos_region` - (Optional) Region of cos bucket. if `type` is `cos`, `cos_region` is required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `code_error` - SCF function code error message.
* `code_result` - SCF function code is correct.
* `code_size` - SCF function code size, unit is M.
* `eip_fixed` - Whether EIP is a fixed IP.
* `eips` - SCF function EIP list.
* `err_no` - SCF function code error code.
* `host` - SCF function domain name.
* `install_dependency` - Whether to automatically install dependencies.
* `modify_time` - SCF function last modified time.
* `status_desc` - SCF status description.
* `status` - SCF function status.
* `trigger_info` - SCF trigger details list. Each element contains the following attributes:
  * `create_time` - Create time of SCF function trigger.
  * `custom_argument` - User-defined parameters of SCF function trigger.
  * `enable` - Whether SCF function trigger is enable.
  * `modify_time` - Modify time of SCF function trigger.
  * `name` - Name of SCF function trigger.
  * `trigger_desc` - TriggerDesc of SCF function trigger.
  * `type` - Type of SCF function trigger.
* `vip` - SCF function vip.


## Import

SCF function can be imported, e.g.

-> **NOTE:** function id is `<function namespace>+<function name>`

```
$ terraform import tencentcloud_scf_function.test default+test
```

