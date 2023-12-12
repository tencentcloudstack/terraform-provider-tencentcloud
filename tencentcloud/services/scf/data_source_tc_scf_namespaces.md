Use this data source to query SCF namespaces.

Example Usage

```hcl
resource "tencentcloud_scf_namespace" "foo" {
  namespace = "ci-test-scf"
}

data "tencentcloud_scf_namespaces" "foo" {
  namespace = tencentcloud_scf_namespace.foo.id
}
```