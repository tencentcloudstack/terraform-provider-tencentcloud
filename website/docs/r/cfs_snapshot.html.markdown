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

* Create: [CreateCfsSnapshot](https://www.tencentcloud.com/document/api/582/50843)
* Create: [DescribeCfsSnapshots](https://www.tencentcloud.com/document/api/582/50836)
* Create: [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Read: [DescribeCfsSnapshots](https://www.tencentcloud.com/document/api/582/50836)
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/api/651/33022)
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/api/651/33012)
* Update: [UpdateCfsSnapshotAttribute](https://www.tencentcloud.com/document/api/582/50832)
* Delete: [DeleteCfsSnapshot](https://www.tencentcloud.com/document/api/582/50842)
* Delete: [DescribeCfsSnapshots](https://www.tencentcloud.com/document/api/582/50836)

## Import

cfs snapshot can be imported using the id, e.g.

```
terraform import tencentcloud_cfs_snapshot.snapshot snapshot_id
```

