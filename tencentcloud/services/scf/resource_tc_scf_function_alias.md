Provides a resource to create a scf function alias

Example Usage

Route by weight

```hcl
resource "tencentcloud_scf_function_alias" "example" {
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
```

Route by matching rules

```hcl
resource "tencentcloud_scf_function_alias" "example" {
  description      = "matchs for test"
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
terraform import tencentcloud_scf_function_alias.example namespace#functionName#name
```