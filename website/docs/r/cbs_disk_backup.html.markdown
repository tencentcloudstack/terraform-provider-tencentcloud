---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_disk_backup"
sidebar_current: "docs-tencentcloud-resource-cbs_disk_backup"
description: |-
  Provides a resource to create a CBS disk backup.
---

# tencentcloud_cbs_disk_backup

Provides a resource to create a CBS disk backup.

~> **NOTE:** The parameter `disk_backup_quota` in the resource `tencentcloud_cbs_storage` must be greater than 1.

## Example Usage

```hcl
resource "tencentcloud_cbs_storage" "example" {
  storage_name      = "tf-example"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-6"
  project_id        = 0
  encrypt           = false
  disk_backup_quota = 3

  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cbs_disk_backup" "example" {
  disk_id          = tencentcloud_cbs_storage.example.id
  disk_backup_name = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `disk_id` - (Required, String, ForceNew) ID of the original cloud disk of the backup point, which can be queried through the DescribeDisks API.
* `disk_backup_name` - (Optional, String, ForceNew) Backup point name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateDiskBackup](https://cloud.tencent.com/document/api/362/84591) - This API is used to create a backup point for a cloud disk.
* Create: [DescribeDiskBackups](https://cloud.tencent.com/document/api/362/80278) - This API is used to query the details of backup points.
* Read: [DescribeDiskBackups](https://cloud.tencent.com/document/api/362/80278) - This API is used to query the details of backup points.
* Delete: [DeleteDiskBackups](https://cloud.tencent.com/document/api/362/80279) - This API is used to delete the backup points of the specified cloud disk in batches.
* Delete: [DescribeDiskBackups](https://cloud.tencent.com/document/api/362/80278) - This API is used to query the details of backup points.

## Import

CBS disk backup can be imported using the id, e.g.

```
terraform import tencentcloud_cbs_disk_backup.example dbp-qax6zwvr
```

