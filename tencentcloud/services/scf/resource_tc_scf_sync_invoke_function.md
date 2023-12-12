Provides a resource to create a scf sync_invoke_function

Example Usage

```hcl
resource "tencentcloud_scf_sync_invoke_function" "invoke_function" {
  function_name = "keep-1676351130"
  qualifier     = "2"
  namespace     = "default"
}
```