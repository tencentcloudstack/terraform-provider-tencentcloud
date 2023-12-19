---
subcategory: "Serverless Cloud Function(SCF)"
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

### Using Zip file

```hcl
resource "tencentcloud_scf_function" "foo" {
  name              = "ci-test-function"
  handler           = "first.do_it_first"
  runtime           = "Python3.6"
  enable_public_net = true
  dns_cache         = true
  intranet_config {
    ip_fixed = "ENABLE"
  }
  vpc_id    = "vpc-391sv4w3"
  subnet_id = "subnet-ljyn7h30"

  zip_file = "/scf/first.zip"

  tags = {
    "env" = "test"
  }
}
```

### Using CFS config

```hcl
resource "tencentcloud_scf_function" "foo" {
  name    = "ci-test-function"
  handler = "main.do_it"
  runtime = "Python3.6"

  cfs_config {
    user_id          = "10000"
    user_group_id    = "10000"
    cfs_id           = "cfs-xxxxxxxx"
    mount_ins_id     = "cfs-xxxxxxxx"
    local_mount_dir  = "/mnt"
    remote_mount_dir = "/"
  }
}
```

### Using triggers

```hcl
resource "tencentcloud_scf_function" "foo" {
  name              = "ci-test-function"
  handler           = "first.do_it_first"
  runtime           = "Python3.6"
  enable_public_net = true

  zip_file = "/scf/first.zip"

  triggers {
    name         = "tf-test-fn-trigger"
    type         = "timer"
    trigger_desc = "*/5 * * * * * *"
  }

  triggers {
    name         = "scf-bucket-1308919341.cos.ap-guangzhou.myqcloud.com"
    cos_region   = "ap-guangzhou"
    type         = "cos"
    trigger_desc = "{\"event\":\"cos:ObjectCreated:Put\",\"filter\":{\"Prefix\":\"\",\"Suffix\":\"\"}}"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Name of the SCF function. Name supports 26 English letters, numbers, connectors, and underscores, it should start with a letter. The last character cannot be `-` or `_`. Available length is 2-60.
* `async_run_enable` - (Optional, String, ForceNew) Whether SCF function asynchronous attribute is enabled. `TRUE` is open, `FALSE` is close.
* `cfs_config` - (Optional, List) List of CFS configurations.
* `cls_logset_id` - (Optional, String) cls logset id of the SCF function.
* `cls_topic_id` - (Optional, String) cls topic id of the SCF function.
* `cos_bucket_name` - (Optional, String) Cos bucket name of the SCF function, such as `cos-1234567890`, conflict with `zip_file`.
* `cos_bucket_region` - (Optional, String) Cos bucket region of the SCF function, conflict with `zip_file`.
* `cos_object_name` - (Optional, String) Cos object name of the SCF function, should have suffix `.zip` or `.jar`, conflict with `zip_file`.
* `description` - (Optional, String) Description of the SCF function. Description supports English letters, numbers, spaces, commas, newlines, periods and Chinese, the maximum length is 1000.
* `dns_cache` - (Optional, Bool) Whether to enable Dns caching capability, only the EVENT function is supported. Default is false.
* `enable_eip_config` - (Optional, Bool) Indicates whether EIP config set to `ENABLE` when `enable_public_net` was true. Default `false`.
* `enable_public_net` - (Optional, Bool) Indicates whether public net config enabled. Default `false`. NOTE: only `vpc_id` specified can disable public net config.
* `environment` - (Optional, Map) Environment of the SCF function.
* `func_type` - (Optional, String) Function type. The default value is Event. Enter Event if you need to create a trigger function. Enter HTTP if you need to create an HTTP function service.
* `handler` - (Optional, String) Handler of the SCF function. The format of name is `<filename>.<method_name>`, and it supports 26 English letters, numbers, connectors, and underscores, it should start with a letter. The last character cannot be `-` or `_`. Available length is 2-60.
* `image_config` - (Optional, List) Image of the SCF function, conflict with `cos_bucket_name`, `cos_object_name`, `cos_bucket_region`, `zip_file`.
* `intranet_config` - (Optional, List) Intranet access configuration.
* `l5_enable` - (Optional, Bool) Enable L5 for SCF function, default is `false`.
* `layers` - (Optional, List) The list of association layers.
* `mem_size` - (Optional, Int) Memory size of the SCF function, unit is MB. The default is `128`MB. The ladder is 128M.
* `namespace` - (Optional, String, ForceNew) Namespace of the SCF function, default is `default`.
* `role` - (Optional, String) Role of the SCF function.
* `runtime` - (Optional, String) Runtime of the SCF function, only supports `Python2.7`, `Python3.6`, `Nodejs6.10`, `Nodejs8.9`, `Nodejs10.15`, `Nodejs12.16`, `Php5.2`, `Php7.4`, `Go1`, `Java8`, and `CustomRuntime`, default is `Python2.7`.
* `subnet_id` - (Optional, String) Subnet ID of the SCF function.
* `tags` - (Optional, Map) Tags of the SCF function.
* `timeout` - (Optional, Int) Timeout of the SCF function, unit is second. Default `3`. Available value is 1-900.
* `triggers` - (Optional, Set) Trigger list of the SCF function, note that if you modify the trigger list, all existing triggers will be deleted, and then create triggers in the new list. Each element contains the following attributes:
* `vpc_id` - (Optional, String) VPC ID of the SCF function.
* `zip_file` - (Optional, String) Zip file of the SCF function, conflict with `cos_bucket_name`, `cos_object_name`, `cos_bucket_region`.

The `cfs_config` object supports the following:

* `cfs_id` - (Required, String) File system instance ID.
* `local_mount_dir` - (Required, String) Local mount directory.
* `mount_ins_id` - (Required, String) File system mount instance ID.
* `remote_mount_dir` - (Required, String) Remote mount directory.
* `user_group_id` - (Required, String) ID of user group.
* `user_id` - (Required, String) ID of user.

The `image_config` object supports the following:

* `image_type` - (Required, String) The image type. personal or enterprise.
* `image_uri` - (Required, String) The uri of image.
* `args` - (Optional, String) the parameters of command.
* `command` - (Optional, String) The command of entrypoint.
* `container_image_accelerate` - (Optional, Bool) Image accelerate switch.
* `entry_point` - (Optional, String) The entrypoint of app.
* `image_port` - (Optional, Int) Image function port setting. Default is `9000`, -1 indicates no port mirroring function. Other value ranges 0 ~ 65535.
* `registry_id` - (Optional, String) The registry id of TCR. When image type is enterprise, it must be set.

The `intranet_config` object supports the following:

* `ip_fixed` - (Required, String) Whether to enable fixed intranet IP, ENABLE is enabled, DISABLE is disabled.

The `layers` object supports the following:

* `layer_name` - (Required, String) The name of Layer.
* `layer_version` - (Required, Int) The version of layer.

The `triggers` object supports the following:

* `name` - (Required, String) Name of the SCF function trigger, if `type` is `ckafka`, the format of name must be `<ckafkaInstanceId>-<topicId>`; if `type` is `cos`, the name is cos bucket id, other In any case, it can be combined arbitrarily. It can only contain English letters, numbers, connectors and underscores. The maximum length is 100.
* `trigger_desc` - (Required, String) TriggerDesc of the SCF function trigger, parameter format of `timer` is linux cron expression; parameter of `cos` type is json string `{"bucketUrl":"<name-appid>.cos.<region>.myqcloud.com","event":"cos:ObjectCreated:*","filter":{"Prefix":"","Suffix":""}}`, where `bucketUrl` is cos bucket (optional), `event` is the cos event trigger, `Prefix` is the corresponding file prefix filter condition, `Suffix` is the suffix filter condition, if not need filter condition can not pass; `cmq` type does not pass this parameter; `ckafka` type parameter format is json string `{"maxMsgNum":"1","offset":"latest"}`; `apigw` type parameter format is json string `{"api":{"authRequired":"FALSE","requestConfig":{"method":"ANY"},"isIntegratedResponse":"FALSE"},"service":{"serviceId":"service-dqzh68sg"},"release":{"environmentName":"test"}}`.
* `type` - (Required, String) Type of the SCF function trigger, support `cos`, `cmq`, `timer`, `ckafka`, `apigw`.
* `cos_region` - (Optional, String) Region of cos bucket. if `type` is `cos`, `cos_region` is required.

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

