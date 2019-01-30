---
layout: "tencentcloud"
page_title: "tencentcloud: tencentcloud_cbs_storage"
sidebar_current: "docs-tencentcloud-resource-cbs-storage-x"
description: |-
  Provides a CBS resource.
---

# tencentcloud_cbs_storage

Provides a CBS resource.

~> **NOTE:** At present, only 'PREPAID' storage is supported to create. 'PREPAID' storage cannot be deleted, once created, must wait it to be expired and release it automatically.

## Example Usage

```hcl
data "tencentcloud_availability_zones" "my_favorate_zones" {}

resource "tencentcloud_cbs_storage" "my-storage" {
  storage_type      = "cloudBasic"
  storage_size      = 50
  period            = 3
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  storage_name      = "my-storage"
}
```

## Argument Reference

The following arguments are supported:

* `storage_type` - (Required) Type of CBS medium. cloudBasic refers to a HDD cloud storage, cloudPremium refers to a Premium cloud storage, cloudSSD refers to a SSD cloud storage. **NOTE**, `storage_type` do not support modification.
* `storage_size` - (Required) Size of the storage (GB). The value range is 10GB - 4,000GB (HDD cloud storages), 500GB - 4,000GB (Premium cloud storages), 100GB - 4,000GB (SSD cloud storages). The increment is 10GB. **NOTE**,  `storage_size` do not support modification.
* `period` - (Required) The tenancy (time unit is month) of the prepaid storage, the legal values are [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60]. **NOTE**, `period` do not support modification.
* `availability_zone` - (Required) The available zone that the CBS instance locates at. **NOTE**, `availability_zone` do not support modification.
* `storage_name` - (Optional) The name of the CBS. This storage_name can have a string of 1 to 64 characters, must contain only alphanumeric characters or hyphens, such as "-",".","_". If not specified, the default name is `CBS-Instance`. It is supported to modify `storage_name` after the storage is created
* `snapshot_id` - (Optional) For a new storage, this indicate which snapshot to use to create the new storage. **For a exist storage, change this field whill case a rollback operation: your storage will rollback to the moment the snapshot created, your must change this filed carefully, please ensure your data in this storage is saved or out of use.**


## Attributes Reference

The following attributes are exported:

* `id` - The storage ID, something looks like `disk-xxxxxx`.
* `storage_type` - Type of CBS medium.
* `storage_size` - Size of the storage.
* `period` - The tenancy of the storage.
* `availability_zone` - The available zone that the CBS instance.
* `storage_status` - The status of storage. The standard values are as follows, normal: Normal, toRecycle: To be terminated, attaching: Mounting, detaching: Unmounting.
* `attached` - The attach status of storage. 1 indicates that storage has been mounted, 0 indicates the storage unmounted.
