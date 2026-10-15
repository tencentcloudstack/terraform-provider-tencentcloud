---
subcategory: "Cloud HDFS(CHDFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_chdfs_mount_point"
sidebar_current: "docs-tencentcloud-resource-chdfs_mount_point"
description: |-
  Provides a resource to create a chdfs mount_point
---

# tencentcloud_chdfs_mount_point

Provides a resource to create a chdfs mount_point

## Example Usage

```hcl
resource "tencentcloud_chdfs_mount_point" "mount_point" {
  file_system_id     = "f14mpfy5lh4e"
  mount_point_name   = "terraform-test"
  mount_point_status = 1
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, String, ForceNew) file system id you want to mount.
* `mount_point_name` - (Required, String) mount point name.
* `mount_point_status` - (Required, Int) mount status 1:open, 2:close.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateMountPoint](https://cloud.tencent.com/document/api/1105/51154) - 创建挂载点
* Read: [DescribeMountPoint](https://cloud.tencent.com/document/api/1105/51153) - 查看挂载点详细信息
* Update: [ModifyMountPoint](https://cloud.tencent.com/document/api/1105/51150) - 修改挂载点属性
* Delete: [DeleteMountPoint](https://cloud.tencent.com/document/api/1105/52342) - 删除挂载点

## Import

chdfs mount_point can be imported using the id, e.g.

```
terraform import tencentcloud_chdfs_mount_point.mount_point mount_point_id
```

