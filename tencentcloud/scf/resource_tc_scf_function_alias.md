Provides a resource to create a scf function_alias

Example Usage

```hcl
// by weight
resource "tencentcloud_scf_function_alias" "function_alias" {
  description      = "weight test"
  function_name    = "keep-1676351130"
  function_version = "$LATEST"
  name             = "weight"
  namespace        = "default"

  routing_config {
    additional_version_weights {
      version = "2"
      weight  = 0.4
    }
  }
}

// by route
resource "tencentcloud_scf_function_alias" "function_alias" {
  description      = "matchs for test 12312312"
  function_name    = "keep-1676351130"
  function_version = "3"
  name             = "matchs"
  namespace        = "default"

  routing_config {
    additional_version_matches {
      expression = "testuser"
      key        = "invoke.headers.User"
      method     = "exact"
      version    = "2"
    }
  }
}
```

Import

scf function_alias can be imported using the id, e.g.

```
terraform import tencentcloud_scf_function_alias.function_alias namespace#functionName#name
```