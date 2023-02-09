---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_disk_backup"
sidebar_current: "docs-tencentcloud-resource-cbs_disk_backup"
description: |-
  Provides a resource to create a cbs disk_backup.
---

# tencentcloud_cbs_disk_backup

Provides a resource to create a cbs disk_backup.

~> **NOTE:** Backup quota must greater than 1.

## Example Usage

```hcl
resource "tencentcloud_cbs_disk_backup" "disk_backup" {
  disk_id          = "disk-xxx"
  disk_backup_name = "xxx"
}
```

## Argument Reference

The following arguments are supported:

* `disk_id` - (Required, String, ForceNew) ID of the original cloud disk of the backup point, which can be queried through the DescribeDisks API.
* `disk_backup_name` - (Optional, String, ForceNew) Backup point name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

cbs disk_backup can be imported using the id, e.g.

```
terraform import tencentcloud_cbs_disk_backup.disk_backup disk_backup_id
```

