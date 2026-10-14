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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ExecuteFunction](https://cloud.tencent.com/document/api/862/48427) - This interface is only used for custom development in special scenarios. Unless a Media Processing Service customer service representative proactively informs you that you need to use this interface, do not call it in other cases.

