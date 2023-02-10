---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_snapshot_share_permission"
sidebar_current: "docs-tencentcloud-resource-cbs_snapshot_share_permission"
description: |-
  Provides a resource to create a cbs snapshot_share_permission
---

# tencentcloud_cbs_snapshot_share_permission

Provides a resource to create a cbs snapshot_share_permission

## Example Usage

```hcl
resource "tencentcloud_cbs_snapshot_share_permission" "snapshot_share_permission" {
  account_ids = ["1xxxxxx", "2xxxxxx"]
  snapshot_id = "snap-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `account_ids` - (Required, Set: [`String`]) List of account IDs with which a snapshot is shared. For the format of array-type parameters, see[API Introduction](https://cloud.tencent.com/document/api/213/568). You can find the account ID in[Account Information](https://console.cloud.tencent.com/developer).
* `snapshot_id` - (Required, String) The ID of the snapshot to be queried. You can obtain this by using [DescribeSnapshots](https://cloud.tencent.com/document/api/362/15647).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cbs snapshot_share_permission can be imported using the id, e.g.

```
terraform import tencentcloud_cbs_snapshot_share_permission.snapshot_share_permission snap-xxxxxx
```

