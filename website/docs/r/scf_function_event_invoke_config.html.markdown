---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_function_event_invoke_config"
sidebar_current: "docs-tencentcloud-resource-scf_function_event_invoke_config"
description: |-
  Provides a resource to create a scf function_event_invoke_config
---

# tencentcloud_scf_function_event_invoke_config

Provides a resource to create a scf function_event_invoke_config

## Example Usage

```hcl
resource "tencentcloud_scf_function_event_invoke_config" "function_event_invoke_config" {
  function_name = "keep-1676351130"
  namespace     = "default"
  async_trigger_config {
    retry_config {
      retry_num = 2
    }
    msg_ttl = 24
  }
}
```

## Argument Reference

The following arguments are supported:

* `async_trigger_config` - (Required, List) Async retry configuration information.
* `function_name` - (Required, String) Function name.
* `namespace` - (Optional, String) Function namespace. Default value: default.

The `async_trigger_config` object supports the following:

* `msg_ttl` - (Required, Int) Message retention period.
* `retry_config` - (Required, List) Async retry configuration of function upon user error.

The `retry_config` object of `async_trigger_config` supports the following:

* `retry_num` - (Required, Int) Number of retry attempts.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

scf function_event_invoke_config can be imported using the id, e.g.

```
terraform import tencentcloud_scf_function_event_invoke_config.function_event_invoke_config function_name#namespace
```

