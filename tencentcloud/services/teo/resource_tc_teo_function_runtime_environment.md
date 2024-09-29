Provides a resource to create a teo teo_function_runtime_environment

Example Usage

```hcl
resource "tencentcloud_teo_function_runtime_environment" "teo_function_runtime_environment" {
    function_id = "ef-txx7fnua"
    zone_id     = "zone-2qtuhspy7cr6"

    environment_variables {
        key   = "test-a"
        type  = "string"
        value = "AAA"
    }
    environment_variables {
        key   = "test-b"
        type  = "string"
        value = "BBB"
    }
}
```

Import

teo teo_function_runtime_environment can be imported using the id, e.g.

```
terraform import tencentcloud_teo_function_runtime_environment.teo_function_runtime_environment zone_id#function_id
```
