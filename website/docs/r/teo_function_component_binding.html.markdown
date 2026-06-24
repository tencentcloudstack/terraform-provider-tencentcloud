---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_function_component_binding"
sidebar_current: "docs-tencentcloud-resource-teo_function_component_binding"
description: |-
  Provides a resource to manage TencentCloud EdgeOne (TEO) function component binding configuration
---

# tencentcloud_teo_function_component_binding

Provides a resource to manage TencentCloud EdgeOne (TEO) function component binding configuration

## Example Usage

```hcl
resource "tencentcloud_teo_function_component_binding" "example" {
  zone_id     = "zone-2q0wu2x2oxxxx"
  function_id = "ef-2q0wu2x2oxxxx"

  function_component_bindings {
    type          = "kv_namespace"
    variable_name = "MY_KV"
    kv_namespace_parameters {
      zone_id   = "zone-2q0wu2x2oxxxx"
      namespace = "my-namespace"
    }
  }

  function_component_bindings {
    type          = "kv_namespace"
    variable_name = "MY_KV_2"
    kv_namespace_parameters {
      zone_id   = "zone-2q0wu2x2oxxxx"
      namespace = "my-namespace-2"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `function_component_bindings` - (Required, List) Function component binding list.
* `function_id` - (Required, String, ForceNew) Function ID.
* `zone_id` - (Required, String, ForceNew) Site ID.

The `function_component_bindings` object supports the following:

* `type` - (Required, String) The type of the bound component. Valid value: `kv_namespace`.
* `variable_name` - (Required, String) The variable name used for binding, limited to 1-50 characters.
* `kv_namespace_parameters` - (Optional, List) KV namespace configuration parameters. Required when type is `kv_namespace`.

The `kv_namespace_parameters` object of `function_component_bindings` supports the following:

* `namespace` - (Required, String) The KV namespace name.
* `zone_id` - (Required, String) The site ID of the KV namespace.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TEO function component binding can be imported using the zoneId#functionId, e.g.

```
terraform import tencentcloud_teo_function_component_binding.example zone-2q0wu2x2oxxxx#ef-2q0wu2x2oxxxx
```

