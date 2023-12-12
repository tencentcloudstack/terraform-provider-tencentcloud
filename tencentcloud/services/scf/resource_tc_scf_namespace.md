Provide a resource to create a SCF namespace.

Example Usage

```hcl
resource "tencentcloud_scf_namespace" "foo" {
  namespace = "ci-test-scf"
}
```

Import

SCF namespace can be imported, e.g.

```
$ terraform import tencentcloud_scf_function.test default
```