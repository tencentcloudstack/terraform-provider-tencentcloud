---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_execute_function_operation"
sidebar_current: "docs-tencentcloud-resource-mps_execute_function_operation"
description: |-
  Provides a resource to create a mps execute_function_operation
---

# tencentcloud_mps_execute_function_operation

Provides a resource to create a mps execute_function_operation

## Example Usage

```hcl
resource "tencentcloud_mps_execute_function_operation" "operation" {
  function_name = "ExampleFunc"
  function_arg  = "arg1"
}
```

## Argument Reference

The following arguments are supported:

* `function_arg` - (Required, String, ForceNew) API parameter. Parameter format will depend on the actual function definition.
* `function_name` - (Required, String, ForceNew) Name of called backend API.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



