Provides a resource to create a TCSS cluster access

Example Usage

```hcl
resource "tencentcloud_tcss_cluster_access" "example" {
  cluster_id = "cls-fdy7hm1q"
  switch_on  = true

  timeouts {
    create = "20m"
    update = "20m"
    delete = "20m"
  }
}
```

Import

TCSS cluster access can be imported using the id, e.g.

```
terraform import tencentcloud_tcss_cluster_access.example cls-fdy7hm1q
```
