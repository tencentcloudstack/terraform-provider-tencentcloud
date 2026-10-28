Provides a resource to manage TEO edge function component binding configuration

~> **NOTE:** If this resource management method is used, it must be the sole method employed; management via other channels—such as the console, API, or CLI—is prohibited. Doing so would result in unintended changes or cause `apply` or `destroy` operations to overwrite existing configurations.

Example Usage

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

Import

teo function_component_binding can be imported using the composite id (zone_id#function_id), e.g.

```
terraform import tencentcloud_teo_function_component_binding.example zone-2qtuhspy7cr6#ef-txx7fnua
```
