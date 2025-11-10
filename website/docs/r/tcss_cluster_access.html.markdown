---
subcategory: "Tencent Container Security Service(TCSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcss_cluster_access"
sidebar_current: "docs-tencentcloud-resource-tcss_cluster_access"
description: |-
  Provides a resource to create a TCSS cluster access
---

# tencentcloud_tcss_cluster_access

Provides a resource to create a TCSS cluster access

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster Id.
* `switch_on` - (Optional, Bool) Whether to enable cluster defend status.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `accessed_status` - Cluster access status.
* `defender_status` - Cluster defender status.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.

## Import

TCSS cluster access can be imported using the id, e.g.

```
terraform import tencentcloud_tcss_cluster_access.example cls-fdy7hm1q
```

