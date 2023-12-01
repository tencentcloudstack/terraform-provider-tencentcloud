---
subcategory: "Cloud HDFS(CHDFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_chdfs_mount_point_attachment"
sidebar_current: "docs-tencentcloud-resource-chdfs_mount_point_attachment"
description: |-
  Provides a resource to create a chdfs mount_point_attachment
---

# tencentcloud_chdfs_mount_point_attachment

Provides a resource to create a chdfs mount_point_attachment

## Example Usage

```hcl
resource "tencentcloud_chdfs_mount_point_attachment" "mount_point_attachment" {
  access_group_ids = [
    "ag-bvmzrbsm",
    "ag-lairqrgr",
  ]
  mount_point_id = "f14mpfy5lh4e-KuiL"
}
```

## Argument Reference

The following arguments are supported:

* `access_group_ids` - (Required, Set: [`String`], ForceNew) associate access group id.
* `mount_point_id` - (Required, String, ForceNew) associate mount point.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

chdfs mount_point_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_mount_point_attachment.mount_point_attachment mount_point_id
```

