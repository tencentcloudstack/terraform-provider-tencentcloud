---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_component_binding"
sidebar_current: "docs-tencentcloud-resource-teo_function_component_binding"
description: |-
  Provides a resource to manage TEO edge function component binding configuration
---

# tencentcloud_teo_function_component_binding

Provides a resource to manage TEO edge function component binding configuration

## Example Usage

```hcl
resource "tencentcloud_teo_function_component_binding" "example" {
  zone_id     = "zone-2qtuhspy7cr6"
  function_id = "ef-txx7fnua"

  function_component_bindings {
    type          = "kv_namespace"
    variable_name = "MY_KV"

    kv_namespace_parameters {
      zone_id   = "zone-2qtuhspy7cr6"
      namespace = "my-namespace"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `function_component_bindings` - (Required, List) List of function component bindings.
* `function_id` - (Required, String, ForceNew) Function ID.
* `zone_id` - (Required, String, ForceNew) Site ID.

The `function_component_bindings` object supports the following:

* `type` - (Required, String) Component type. Valid values: `kv_namespace`.
* `variable_name` - (Required, String) Variable name for binding, 1-50 characters, alphanumeric and underscore, cannot start with a number.
* `kv_namespace_parameters` - (Optional, List) KV namespace configuration parameters. Required when type is `kv_namespace`.

The `kv_namespace_parameters` object of `function_component_bindings` supports the following:

* `namespace` - (Required, String) KV namespace name.
* `zone_id` - (Required, String) Zone ID of the KV namespace.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

teo function_component_binding can be imported using the composite id (zone_id#function_id), e.g.

```
terraform import tencentcloud_teo_function_component_binding.example zone-2qtuhspy7cr6#ef-txx7fnua
```

