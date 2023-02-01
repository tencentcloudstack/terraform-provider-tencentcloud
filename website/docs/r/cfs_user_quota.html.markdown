---
subcategory: "Cloud File Storage(CFS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cfs_user_quota"
sidebar_current: "docs-tencentcloud-resource-cfs_user_quota"
description: |-
  Provides a resource to create a cfs user_quota
---

# tencentcloud_cfs_user_quota

Provides a resource to create a cfs user_quota

## Example Usage

```hcl
resource "tencentcloud_cfs_user_quota" "user_quota" {
  file_system_id      = "cfs-4636029bc"
  user_type           = "Uid"
  user_id             = "2159973417"
  capacity_hard_limit = 10
  file_hard_limit     = 10000
}
```

## Argument Reference

The following arguments are supported:

* `file_system_id` - (Required, String) File system ID.
* `user_id` - (Required, String) Info of UID/GID.
* `user_type` - (Required, String) Quota type. Valid value: `Uid`, `Gid`.
* `capacity_hard_limit` - (Optional, Int) Capacity Limit(GB).
* `file_hard_limit` - (Optional, Int) File limit.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cfs user_quota can be imported using the id, e.g.

```
terraform import tencentcloud_cfs_user_quota.user_quota user_quota_id
```

