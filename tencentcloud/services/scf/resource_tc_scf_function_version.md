Provides a resource to create a scf function_version

Example Usage

```hcl
resource "tencentcloud_scf_function_version" "function_version" {
  function_name    = "keep-1676351130"
  namespace        = "default"
  description      = "for-terraform-test"
}

```

Import

scf function_version can be imported using the id, e.g.

```
terraform import tencentcloud_scf_function_version.function_version functionName#namespace#functionVersion
```