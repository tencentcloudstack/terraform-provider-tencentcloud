---
subcategory: "Serverless Cloud Function(SCF)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_scf_invoke_function"
sidebar_current: "docs-tencentcloud-resource-scf_invoke_function"
description: |-
  Provides a resource to create a scf invoke_function
---

# tencentcloud_scf_invoke_function

Provides a resource to create a scf invoke_function

## Example Usage

```hcl
resource "tencentcloud_scf_invoke_function" "invoke_function" {
  function_name = "keep-1676351130"
  qualifier     = "2"
  namespace     = "default"
}
```

## Argument Reference

The following arguments are supported:

* `function_name` - (Required, String, ForceNew) Function name.
* `client_context` - (Optional, String, ForceNew) Function running parameter, which is in the JSON format. The maximum parameter size is 6 MB for synchronized invocations and 128KB for asynchronized invocations. This field corresponds to event input parameter.
* `invocation_type` - (Optional, String, ForceNew) Fill in RequestResponse for synchronized invocations (default and recommended) and Event for asychronized invocations. Note that for synchronized invocations, the max timeout period is 300s. Choose asychronized invocations if the required timeout period is longer than 300 seconds. You can also use InvokeFunction for synchronized invocations.
* `log_type` - (Optional, String, ForceNew) Null for async invocations.
* `namespace` - (Optional, String, ForceNew) Namespace.
* `qualifier` - (Optional, String, ForceNew) The version or alias of the triggered function. It defaults to $LATEST.
* `routing_key` - (Optional, String, ForceNew) Traffic routing config in json format, e.g., {k:v}. Please note that both k and v must be strings. Up to 1024 bytes allowed.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



