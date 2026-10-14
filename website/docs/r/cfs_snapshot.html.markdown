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

* Create: [CreateCfsSnapshot](https://cloud.tencent.com/document/api/582/80200) - This API is used to create a file system snapshot.
* Create: [DescribeCfsSnapshots](https://cloud.tencent.com/document/api/582/80206) - This API is used to query the list of snapshots of a file system.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Read: [DescribeCfsSnapshots](https://cloud.tencent.com/document/api/582/80206) - This API is used to query the list of snapshots of a file system.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Update: [UpdateCfsSnapshotAttribute](https://cloud.tencent.com/document/api/582/80202) - This API is used to update the name and retention period of a file system snapshot.
* Delete: [DeleteCfsSnapshot](https://cloud.tencent.com/document/api/582/80199) - This API is used to delete a file system snapshot.
* Delete: [DescribeCfsSnapshots](https://cloud.tencent.com/document/api/582/80206) - This API is used to query the list of snapshots of a file system.

## Import

cfs snapshot can be imported using the id, e.g.

```
terraform import tencentcloud_cfs_snapshot.snapshot snapshot_id
```

