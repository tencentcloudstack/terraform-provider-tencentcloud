---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_disk_backup_rollback_operation"
sidebar_current: "docs-tencentcloud-resource-cbs_disk_backup_rollback_operation"
description: |-
  Provides a resource to rollback cbs disk backup.
---

# tencentcloud_cbs_disk_backup_rollback_operation

Provides a resource to rollback cbs disk backup.

## Example Usage

```hcl
resource "tencentcloud_cbs_disk_backup_rollback_operation" "operation" {
  disk_backup_id = "dbp-xxx"
  disk_id        = "disk-xxx"
}
```

## Argument Reference

The following arguments are supported:

* `disk_backup_id` - (Required, String, ForceNew) Cloud disk backup point ID.
* `disk_id` - (Required, String, ForceNew) Cloud disk backup point original cloud disk ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `is_rollback_completed` - Whether the rollback is completed. `true` meaing rollback completed, `false` meaning still rollbacking.


