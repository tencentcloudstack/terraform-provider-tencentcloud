---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_disk"
sidebar_current: "docs-tencentcloud-resource-lighthouse_disk"
description: |-
  Provides a resource to create a lighthouse disk
---

# tencentcloud_lighthouse_disk

Provides a resource to create a lighthouse disk

## Example Usage

```hcl
resource "tencentcloud_lighthouse_disk" "disk" {
  zone      = "ap-hongkong-2"
  disk_size = 20
  disk_type = "CLOUD_SSD"
  disk_charge_prepaid {
    period     = 1
    renew_flag = "NOTIFY_AND_AUTO_RENEW"
    time_unit  = "m"

  }
  disk_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `disk_charge_prepaid` - (Required, List) Disk subscription related parameter settings.
* `disk_size` - (Required, Int) Disk size, unit: GB.
* `disk_type` - (Required, String) Disk type. Value:CLOUD_PREMIUM, CLOUD_SSD.
* `zone` - (Required, String) Availability zone.
* `auto_mount_configuration` - (Optional, List) Automatically mount and initialize data disks.
* `auto_voucher` - (Optional, Bool) Whether to automatically use the voucher. Not used by default.
* `disk_backup_quota` - (Optional, Int) Specify the disk backup quota. If not uploaded, the default is no backup quota. Currently, only one disk backup quota is supported.
* `disk_count` - (Optional, Int) Disk count. Values: [1, 30]. Default: 1.
* `disk_name` - (Optional, String) Disk name. Maximum length 60.

The `auto_mount_configuration` object supports the following:

* `instance_id` - (Required, String) Instance ID to be mounted. The specified instance must be in the Running state.
* `file_system_type` - (Optional, String) The file system type. Value: ext4, xfs. Only instances of the Linux operating system can pass in this parameter, and if it is not passed, it defaults to ext4.
* `mount_point` - (Optional, String) The mount point within the instance. Only instances of the Linux operating system can pass in this parameter, and if it is not passed, it will be mounted under the /data/disk path by default.

The `disk_charge_prepaid` object supports the following:

* `period` - (Required, Int) new purchase cycle.
* `renew_flag` - (Optional, String) Automatic renewal flag. Value: `NOTIFY_AND_AUTO_RENEW`: Notice expires and auto-renews. `NOTIFY_AND_MANUAL_RENEW`: Notification expires without automatic renewal, users need to manually renew. `DISABLE_NOTIFY_AND_AUTO_RENEW`: No automatic renewal and no notification. Default: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the disk will be automatically renewed monthly when the account balance is sufficient.
* `time_unit` - (Optional, String) newly purchased unit. Default: m.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateDisks](https://www.tencentcloud.com/search/?wd=CreateDisks&type=document&lang=en)
* Create: [DescribeDisks](https://www.tencentcloud.com/search/?wd=DescribeDisks&type=document&lang=en)
* Read: [DescribeDisks](https://www.tencentcloud.com/search/?wd=DescribeDisks&type=document&lang=en)
* Update: [ModifyDisksAttribute](https://www.tencentcloud.com/search/?wd=ModifyDisksAttribute&type=document&lang=en)
* Update: [ModifyDisksRenewFlag](https://www.tencentcloud.com/search/?wd=ModifyDisksRenewFlag&type=document&lang=en)
* Delete: [DescribeDisks](https://www.tencentcloud.com/search/?wd=DescribeDisks&type=document&lang=en)
* Delete: [IsolateDisks](https://www.tencentcloud.com/search/?wd=IsolateDisks&type=document&lang=en)
* Delete: [TerminateDisks](https://www.tencentcloud.com/search/?wd=TerminateDisks&type=document&lang=en)

