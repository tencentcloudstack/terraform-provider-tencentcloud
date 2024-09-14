---
subcategory: "Tencent Kubernetes Engine(TKE)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_kubernetes_health_check_policy"
sidebar_current: "docs-tencentcloud-resource-kubernetes_health_check_policy"
description: |-
  Provides a resource to create a tke kubernetes_health_check_policy
---

# tencentcloud_kubernetes_health_check_policy

Provides a resource to create a tke kubernetes_health_check_policy

## Example Usage

```hcl
resource "tencentcloud_kubernetes_health_check_policy" "kubernetes_health_check_policy" {
  cluster_id = "cls-xxxxx"
  name       = "example"
  rules {
    name                = "OOMKilling"
    auto_repair_enabled = true
    enabled             = true
  }
  rules {
    name                = "KubeletUnhealthy"
    auto_repair_enabled = true
    enabled             = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) ID of the cluster.
* `name` - (Required, String) Health Check Policy Name.
* `rules` - (Required, List) Health check policy rule list.

The `rules` object supports the following:

* `auto_repair_enabled` - (Required, Bool) Enable repair or not.
* `enabled` - (Required, Bool) Enable detection of this project or not.
* `name` - (Required, String) Health check rule details.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tke kubernetes_health_check_policy can be imported using the id, e.g.

```
terraform import tencentcloud_kubernetes_health_check_policy.kubernetes_health_check_policy cls-xxxxx#healthcheckpolicyname
```

