---
subcategory: "Cloud HDFS(CHDFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_chdfs_mount_points"
sidebar_current: "docs-tencentcloud-datasource-chdfs_mount_points"
description: |-
  Use this data source to query detailed information of chdfs mount_points
---

# tencentcloud_chdfs_mount_points

Use this data source to query detailed information of chdfs mount_points

## Example Usage

```hcl
data "tencentcloud_chdfs_mount_points" "mount_points" {
  file_system_id = "f14mpfy5lh4e"
}
```

## Argument Reference

The following arguments are supported:

* `access_group_id` - (Optional, String) get mount points belongs to access group id, only can use one of the AccessGroupId,FileSystemId,OwnerUin paramaters.
* `file_system_id` - (Optional, String) get mount points belongs to file system id, only can use one of the AccessGroupId,FileSystemId,OwnerUin paramaters.
* `owner_uin` - (Optional, Int) get mount points belongs to owner uin, only can use one of the AccessGroupId,FileSystemId,OwnerUin paramaters.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `mount_points` - mount point list.
  * `access_group_ids` - associated group ids.
  * `create_time` - create time.
  * `file_system_id` - mounted file system id.
  * `mount_point_id` - mount point id.
  * `mount_point_name` - mount point name.
  * `status` - mount point status.


