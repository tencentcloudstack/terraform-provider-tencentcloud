---
subcategory: "Cloud File Storage(CFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfs_auto_snapshot_policy_attachment"
sidebar_current: "docs-tencentcloud-resource-cfs_auto_snapshot_policy_attachment"
description: |-
  Provides a resource to create a cfs auto_snapshot_policy_attachment
---

# tencentcloud_cfs_auto_snapshot_policy_attachment

Provides a resource to create a cfs auto_snapshot_policy_attachment

## Example Usage

```hcl
resource "tencentcloud_cfs_auto_snapshot_policy_attachment" "auto_snapshot_policy_attachment" {
  auto_snapshot_policy_id = "asp-basic"
  file_system_ids         = "cfs-4xzkct19,cfs-iobiaxtj"
}
```

## Argument Reference

The following arguments are supported:

* `auto_snapshot_policy_id` - (Required, String, ForceNew) ID of the snapshot to be unbound.
* `file_system_ids` - (Required, String, ForceNew) List of IDs of the file systems to be unbound, separated by comma.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [BindAutoSnapshotPolicy](https://www.tencentcloud.com/document/product/582/50841) - This API is used to bind one or multiple file systems to a snapshot policy. A file system can be bound to only one policy.
* Read: [DescribeAutoSnapshotPolicies](https://www.tencentcloud.com/document/product/582/50838) - This API is used to query the list of scheduled snapshot policies of a file system.
* Delete: [UnbindAutoSnapshotPolicy](https://www.tencentcloud.com/document/product/582/50834) - This API is used to unbind a snapshot policy from a file system.

## Import

cfs auto_snapshot_policy_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_cfs_auto_snapshot_policy_attachment.auto_snapshot_policy_attachment auto_snapshot_policy_id#file_system_ids
```

