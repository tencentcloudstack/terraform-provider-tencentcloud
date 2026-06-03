Provides a resource to manage TEO edge function component binding configuration.

Example Usage

Bindkv namespace to edge function

```hcl
resource "tencentcloud_teo_function_component_binding" "example" {
  zone_id     = "zone-2qtuhspy7cr6"
  function_id = "ef-2qt00hbm7crb"

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

Bind multiple kv namespaces to edge function

```hcl
resource "tencentcloud_teo_function_component_binding" "example" {
  zone_id     = "zone-2qtuhspy7cr6"
  function_id = "ef-2qt00hbm7crb"

  function_component_bindings {
    type          = "kv_namespace"
    variable_name = "MY_KV_1"

    kv_namespace_parameters {
      zone_id   = "zone-2qtuhspy7cr6"
      namespace = "namespace-1"
    }
  }

  function_component_bindings {
    type          = "kv_namespace"
    variable_name = "MY_KV_2"

    kv_namespace_parameters {
      zone_id   = "zone-2qtuhspy7cr6"
      namespace = "namespace-2"
    }
  }
}
```

Import

TEO function component binding config can be imported using the composite ID format `zone_id#function_id`, e.g.

```
terraform import tencentcloud_teo_function_component_binding.example zone-2qtuhspy7cr6#ef-2qt00hbm7crb
```
