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

## Example Usage

```hcl
resource "tencentcloud_scf_function" "function" {
  name            = "tf-example-function"
  runtime         = "Nodejs16.13"
  handler         = "index.main"
  memory_size     = 128
  timeout         = 3
  cos_bucket_name = "tf-example-bucket-1300000000"
  cos_object_name = "function.zip"
}

resource "tencentcloud_scf_trigger" "timer" {
  function_name   = tencentcloud_scf_function.function.name
  namespace       = "default"
  trigger_name    = "tf-example-trigger"
  type            = "timer"
  trigger_desc    = jsonencode({ cron = "*/5 * * * * * *" })
  enable          = "OPEN"
  description     = "tf example trigger"
  qualifier       = "$DEFAULT"
  custom_argument = "Information"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String, ForceNew) Name of the SCF function that the trigger binds to.
* `trigger_name` - (Required, String, ForceNew) Name of the trigger.
* `type` - (Required, String, ForceNew) Trigger type. Valid values: `cos`, `cls`, `timer`, `ckafka`, `http`.
* `custom_argument` - (Optional, String) User custom parameter, only supported by timer trigger.
* `description` - (Optional, String) Trigger description.
* `enable` - (Optional, String) Trigger enable status. Valid values: `OPEN` (enabled), `CLOSE` (disabled).
* `namespace` - (Optional, String, ForceNew) Function namespace. Defaults to `default`.
* `qualifier` - (Optional, String) Function version or alias that the trigger points to. Defaults to `$LATEST`.
* `trigger_desc` - (Optional, String) Trigger description parameter, see the trigger description documentation for details.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `add_time` - Trigger creation time.
* `available_status` - Trigger available status.
* `mod_time` - Trigger last modified time.


## Import

SCF trigger can be imported using the composite id `function_name#namespace#trigger_name`, e.g.

```
terraform import tencentcloud_scf_trigger.timer tf-example-function#default#tf-example-trigger
```

