---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_apply_disk_backup"
sidebar_current: "docs-tencentcloud-resource-lighthouse_apply_disk_backup"
description: |-
  Provides a resource to create a lighthouse apply_disk_backup
---

# tencentcloud_lighthouse_apply_disk_backup

Provides a resource to create a lighthouse apply_disk_backup

## Example Usage

```hcl
resource "tencentcloud_lighthouse_apply_disk_backup" "apply_disk_backup" {
  disk_id        = "lhdisk-xxxxxx"
  disk_backup_id = "lhbak-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `disk_backup_id` - (Required, String, ForceNew) Disk backup ID.
* `disk_id` - (Required, String, ForceNew) Disk ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



