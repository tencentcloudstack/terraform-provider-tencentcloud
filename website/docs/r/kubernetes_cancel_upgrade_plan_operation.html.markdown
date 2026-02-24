---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_cancel_upgrade_plan_operation"
sidebar_current: "docs-tencentcloud-resource-kubernetes_cancel_upgrade_plan_operation"
description: |-
  Provides a resource to create a TKE kubernetes cancel upgrade plan operation
---

# tencentcloud_kubernetes_cancel_upgrade_plan_operation

Provides a resource to create a TKE kubernetes cancel upgrade plan operation

## Example Usage

```hcl
resource "tencentcloud_kubernetes_cancel_upgrade_plan_operation" "example" {
  cluster_id = "cls-d2cit6no"
  plan_id    = 39
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID.
* `plan_id` - (Required, Int, ForceNew) Upgrade plan ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



