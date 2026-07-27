---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_trigger"
sidebar_current: "docs-tencentcloud-resource-scf_trigger"
description: |-
  Provides a resource to create a SCF (Serverless Cloud Function) trigger.
---

# tencentcloud_scf_trigger

Provides a resource to create a SCF (Serverless Cloud Function) trigger.

~> **NOTE:** Using it alongside `tencentcloud_scf_function.trigger` will cause a conflict.

~> **NOTE:** Using it alongside `tencentcloud_scf_trigger_config` will cause a conflict.

## Example Usage

### Create an HTTP-type trigger

```hcl
resource "tencentcloud_scf_function" "foo" {
  async_run_enable  = "FALSE"
  cos_bucket_name   = "tf-test-1308726196"
  cos_bucket_region = "ap-guangzhou"
  cos_object_name   = "/tf-test.zip"
  description       = "helloworld-tftest"
  dns_cache         = false
  enable_eip_config = false
  enable_public_net = true
  environment       = {}
  func_type         = "Event"
  handler           = "index.main_handler"
  l5_enable         = false
  mem_size          = 128
  name              = "helloworld-tf"
  namespace         = "default"
  role              = null
  runtime           = "Python3.9"
  subnet_id         = null
  tags              = {}
  timeout           = 3
  vpc_id            = null
  zip_file          = null
  intranet_config {
    ip_fixed = "DISABLE"
  }
}

resource "tencentcloud_scf_function_version" "function_version" {
  function_name = tencentcloud_scf_function.foo.name
  namespace     = tencentcloud_scf_function.foo.namespace
  description   = "tftest1"
}

resource "tencentcloud_scf_trigger" "http" {
  enable        = "OPEN"
  function_name = tencentcloud_scf_function.foo.name
  namespace     = tencentcloud_scf_function.foo.namespace
  qualifier     = tencentcloud_scf_function_version.function_version.function_version
  trigger_desc = jsonencode({
    ApiGwCompatible = true
    AuthType        = "NONE"
    CorsConfig = {
      Credentials = false
      Enable      = false
      MaxAge      = 0
    }
    EnableSimpleMode = false
    NetConfig = {
      EnableExtranet = false
      EnableIntranet = true
    }
  })
  type = "http"
}
```

### Create a timer-type trigger.

```hcl
resource "tencentcloud_scf_function" "foo" {
  async_run_enable  = "FALSE"
  cos_bucket_name   = "tf-test-1308726196"
  cos_bucket_region = "ap-guangzhou"
  cos_object_name   = "/tf-test.zip"
  description       = "helloworld-tftest"
  dns_cache         = false
  enable_eip_config = false
  enable_public_net = true
  environment       = {}
  func_type         = "Event"
  handler           = "index.main_handler"
  l5_enable         = false
  mem_size          = 128
  name              = "helloworld-tf"
  namespace         = "default"
  role              = null
  runtime           = "Python3.9"
  subnet_id         = null
  tags              = {}
  timeout           = 3
  vpc_id            = null
  zip_file          = null
  intranet_config {
    ip_fixed = "DISABLE"
  }
}

resource "tencentcloud_scf_function_version" "function_version" {
  function_name = tencentcloud_scf_function.foo.name
  namespace     = tencentcloud_scf_function.foo.namespace
  description   = "tftest1"
}

resource "tencentcloud_scf_trigger" "timer" {
  function_name   = tencentcloud_scf_function.foo.name
  namespace       = tencentcloud_scf_function.foo.namespace
  qualifier       = tencentcloud_scf_function_version.function_version.function_version
  custom_argument = null
  description     = null
  enable          = "OPEN"
  trigger_desc    = "0 0 0 */1 * * *"
  trigger_name    = "tf-test-timer"
  type            = "timer"
}
```

### Create a ckafka-type trigger.

```hcl
resource "tencentcloud_scf_function" "foo" {
  async_run_enable  = "FALSE"
  cos_bucket_name   = "tf-test-1308726196"
  cos_bucket_region = "ap-guangzhou"
  cos_object_name   = "/tf-test.zip"
  description       = "helloworld-tftest"
  dns_cache         = false
  enable_eip_config = false
  enable_public_net = true
  environment       = {}
  func_type         = "Event"
  handler           = "index.main_handler"
  l5_enable         = false
  mem_size          = 128
  name              = "helloworld-tf"
  namespace         = "default"
  role              = null
  runtime           = "Python3.9"
  subnet_id         = null
  tags              = {}
  timeout           = 3
  vpc_id            = null
  zip_file          = null
  intranet_config {
    ip_fixed = "DISABLE"
  }
}

resource "tencentcloud_scf_function_version" "function_version" {
  function_name = tencentcloud_scf_function.foo.name
  namespace     = tencentcloud_scf_function.foo.namespace
  description   = "tftest1"
}
resource "tencentcloud_scf_trigger" "kafka" {
  custom_argument = null
  description     = "tf test"
  enable          = "OPEN"
  function_name   = tencentcloud_scf_function.foo.name
  namespace       = tencentcloud_scf_function.foo.namespace
  qualifier       = tencentcloud_scf_function_version.function_version.function_version
  trigger_desc = jsonencode({
    VIPList           = ["11.135.14.109:26246"]
    consumerGroupName = "scf_ckafka-4kdjmwvztf-testhelloworld-tf_DEFAULT"
    instanceId        = "ckafka-4kdjmwvz"
    isInSequence      = "Yes"
    kafkaStatus       = "NORMAL"
    maxMsgNum         = 100
    maxSingleMsgSize  = 100
    maxTotalMsgSize   = 100
    offset            = "latest"
    partitionNum      = 3
    retry             = 10000
    timeOut           = 50
    topicId           = "topic-bjmcknbs"
    topicName         = "tf-test"
    topicStatus       = "NORMAL"
  })
  trigger_name = "ckafka-4kdjmwvz-tf-test"
  type         = "ckafka"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String, ForceNew) Name of the SCF function that the trigger binds to.
* `namespace` - (Required, String, ForceNew) Function namespace. Defaults to `default`.
* `type` - (Required, String, ForceNew) Trigger type. Valid values: `cos`, `timer`, `ckafka`, `http`.To create Function URL please refer to [Creating a Function URL](https://www.tencentcloud.com/document/product/583/69492?lang=en&pg=); To create a CLS trigger, please refer to [Create Deliver CloudFunction (SCF)](https://www.tencentcloud.com/zh/document/product/614/59903).
* `custom_argument` - (Optional, String) User custom parameter, only supported by timer trigger.
* `description` - (Optional, String) Trigger description.
* `enable` - (Optional, String) Trigger enable status. Valid values: `OPEN` (enabled), `CLOSE` (disabled).
* `qualifier` - (Optional, String) Function version or alias that the trigger points to. Defaults to `$LATEST`.
* `trigger_desc` - (Optional, String) Trigger description parameter, see the trigger description documentation for details: https://www.tencentcloud.com/document/product/583/34880.
* `trigger_name` - (Optional, String) Name of the trigger. It must not be specified when `type` is `http`; it must be specified when `type` is not `http`; when `type` is `cos`, it must match `<bucket>.cos.<region>.myqcloud.com`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `add_time` - Trigger creation time.
* `available_status` - Trigger available status.
* `mod_time` - Trigger last modified time.
* `trigger_attribute` - Trigger description parameter, see the trigger description documentation for details: https://www.tencentcloud.com/document/product/583/34880.


## Import

SCF trigger can be imported using the composite id `function_name#namespace#trigger_name`, e.g.

```
terraform import tencentcloud_scf_trigger.timer tf-example-function#default#tf-example-trigger
```

Or use Terraform 1.5+ `import` block:

```hcl
import {
  to = tencentcloud_scf_trigger.http
  id = "tf-example-function#default#tf-example-trigger"
}
```

