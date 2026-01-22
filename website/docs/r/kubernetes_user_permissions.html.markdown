---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_user_permissions"
sidebar_current: "docs-tencentcloud-resource-kubernetes_user_permissions"
description: |-
  Provides a resource to create a TKE kubernetes user permissions
---

# tencentcloud_kubernetes_user_permissions

Provides a resource to create a TKE kubernetes user permissions

~> **NOTE:** This resource must exclusive in one target Uin, do not declare additional permissions resources of this target Uin elsewhere.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `target_uin` - (Required, String, ForceNew) Unique identifier of the user to be authorized (supports sub-account UIN and role UIN).
* `permissions` - (Optional, Set) Complete list of permissions that the user should ultimately have. Uses declarative semantics, the passed list represents all permissions the user should ultimately have, the system will automatically calculate differences and perform necessary create/delete operations. When empty or not provided, all permissions for this user will be cleared. Maximum support for 100 permission items.

The `permissions` object supports the following:

* `cluster_id` - (Required, String) Cluster ID.
* `role_name` - (Required, String) Role name. Predefined roles include: tke:admin (cluster administrator), tke:ops (operations personnel), tke:dev (developer), tke:ro (read-only user), tke:ns:dev (namespace developer), tke:ns:ro (namespace read-only user), others are user-defined roles.
* `role_type` - (Required, String) Authorization type. Enum values: cluster (cluster-level permissions, corresponding to ClusterRoleBinding), namespace (namespace-level permissions, corresponding to RoleBinding).
* `is_custom` - (Optional, Bool) Whether it is a custom role, default false.
* `namespace` - (Optional, String) Namespace. Required when RoleType is namespace.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

TKE kubernetes user permissions can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_user_permissions.example 100056451191
```

