---
subcategory: "Cloud File Storage(CFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfs_snapshot"
sidebar_current: "docs-tencentcloud-resource-cfs_snapshot"
description: |-
  Provides a resource to create a cfs snapshot
---

# tencentcloud_cfs_snapshot

Provides a resource to create a cfs snapshot

## Example Usage

```hcl
resource "tencentcloud_cfs_snapshot" "snapshot" {
  file_system_id = "cfs-iobiaxtj"
  snapshot_name  = "test"
  tags = {
    "createdBy" = "terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, String, ForceNew) Id of file system.
* `snapshot_name` - (Optional, String) Name of snapshot.
* `tags` - (Optional, Map) Tag description list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateCfsSnapshot](https://cloud.tencent.com/document/api/582/80200) - 创建文件系统快照
* Create: [DescribeCfsSnapshots](https://cloud.tencent.com/document/api/582/80206) - 查询快照列表
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Read: [DescribeCfsSnapshots](https://cloud.tencent.com/document/api/582/80206) - 查询快照列表
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - 批量查看资源关联的标签
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Update: [UpdateCfsSnapshotAttribute](https://cloud.tencent.com/document/api/582/80202) - 更新文件系统快照信息
* Delete: [DeleteCfsSnapshot](https://cloud.tencent.com/document/api/582/80199) - 删除文件系统快照
* Delete: [DescribeCfsSnapshots](https://cloud.tencent.com/document/api/582/80206) - 查询快照列表

## Import

cfs snapshot can be imported using the id, e.g.

```
terraform import tencentcloud_cfs_snapshot.snapshot snapshot_id
```

