---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_disk_backup"
sidebar_current: "docs-tencentcloud-resource-lighthouse_disk_backup"
description: |-
  Provides a resource to create a lighthouse disk_backup
---

# tencentcloud_lighthouse_disk_backup

Provides a resource to create a lighthouse disk_backup

## Example Usage

```hcl
resource "tencentcloud_lighthouse_disk_backup" "disk_backup" {
  disk_id          = "lhdisk-xxxxx"
  disk_backup_name = "disk-backup"
}
```

## Argument Reference

The following arguments are supported:

* `disk_id` - (Required, String, ForceNew) Disk ID. Only data disks are supported to create disk backup.
* `disk_backup_name` - (Optional, String) Disk backup name. The maximum length is 90 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateDiskBackup](https://cloud.tencent.com/document/api/1207/84381)
* Create: [DescribeDiskBackups](https://cloud.tencent.com/document/api/1207/84379)
* Read: [DescribeDiskBackups](https://cloud.tencent.com/document/api/1207/84379)
* Update: [ModifyDiskBackupsAttribute](https://cloud.tencent.com/document/api/1207/84377)
* Delete: [DeleteDiskBackups](https://cloud.tencent.com/document/api/1207/84380)

## Import

lighthouse disk_backup can be imported using the id, e.g.

```
terraform import tencentcloud_lighthouse_disk_backup.disk_backup disk_backup_id
```

