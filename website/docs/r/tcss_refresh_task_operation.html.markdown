---
subcategory: "Tencent Container Security Service(TCSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcss_refresh_task_operation"
sidebar_current: "docs-tencentcloud-resource-tcss_refresh_task_operation"
description: |-
  Provides a resource to create a TCSS refresh task operation
---

# tencentcloud_tcss_refresh_task_operation

Provides a resource to create a TCSS refresh task operation

## Example Usage

```hcl
resource "tencentcloud_tcss_refresh_task_operation" "example" {}
```

### Or

```hcl
resource "tencentcloud_tcss_refresh_task_operation" "example" {
  cluster_ids = [
    "cls-fdy7hm1q"
  ]
  is_sync_list_only = false
}
```

## Argument Reference

The following arguments are supported:

* `cluster_ids` - (Optional, Set: [`String`], ForceNew) Cluster Id list.
* `is_sync_list_only` - (Optional, Bool, ForceNew) Whether to sync list only.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



