Provides a resource to manage TencentCloud EdgeOne (TEO) function component binding configuration

Example Usage

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

Import

TEO function component binding can be imported using the zoneId#functionId, e.g.

```
terraform import tencentcloud_teo_function_component_binding.example zone-2q0wu2x2oxxxx#ef-2q0wu2x2oxxxx
```
