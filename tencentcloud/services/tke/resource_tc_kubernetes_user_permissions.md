Provides a resource to create a TKE kubernetes user permissions

~> **NOTE:** This resource must exclusive in one target Uin, do not declare additional permissions resources of this target Uin elsewhere.

Example Usage

```hcl
resource "tencentcloud_kubernetes_user_permissions" "example" {
  target_uin = "100056451191"
  permissions {
    cluster_id = "cls-62ch3v24"
    role_name  = "tke:admin"
    role_type  = "cluster"
    is_custom  = false
  }

  permissions {
    cluster_id = "cls-62ch3v24"
    role_name  = "tke:admin"
    role_type  = "namespace"
    is_custom  = false
    namespace  = "default"
  }
}
```

Import

TKE kubernetes user permissions can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_user_permissions.example 100056451191
```
