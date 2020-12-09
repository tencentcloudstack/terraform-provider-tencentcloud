---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_storage"
sidebar_current: "docs-tencentcloud-resource-cbs_storage"
description: |-
  Provides a resource to create a CBS.
---

# tencentcloud_cbs_storage

Provides a resource to create a CBS.

## Example Usage

```hcl
resource "tencentcloud_cbs_storage" "storage" {
  storage_name      = "mystorage"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = false

  tags = {
    test = "tf"
  }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The available zone that the CBS instance locates at.
* `storage_name` - (Required) Name of CBS. The maximum length can not exceed 60 bytes.
* `storage_size` - (Required) Volume of CBS, and unit is GB. If storage type is `CLOUD_SSD`, the size range is [100, 16000], and the others are [10-16000].
* `storage_type` - (Required, ForceNew) Type of CBS medium. Valid values: CLOUD_BASIC, CLOUD_PREMIUM and CLOUD_SSD.
* `charge_type` - (Optional) The charge type of CBS instance. Valid values are `PREPAID` and `POSTPAID_BY_HOUR`. The default is `POSTPAID_BY_HOUR`.
* `encrypt` - (Optional, ForceNew) Indicates whether CBS is encrypted.
* `force_delete` - (Optional) Indicate whether to delete CBS instance directly or not. Default is false. If set true, the instance will be deleted instead of staying recycle bin.
* `period` - (Optional, **Deprecated**) It has been deprecated from version 1.33.0. Set `prepaid_period` instead. The purchased usage period of CBS. Valid values: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36].
* `prepaid_period` - (Optional) The tenancy (time unit is month) of the prepaid instance, NOTE: it only works when charge_type is set to `PREPAID`. Valid values are 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36.
* `prepaid_renew_flag` - (Optional) Auto Renewal flag. Value range: `NOTIFY_AND_AUTO_RENEW`: Notify expiry and renew automatically, `NOTIFY_AND_MANUAL_RENEW`: Notify expiry but do not renew automatically, `DISABLE_NOTIFY_AND_MANUAL_RENEW`: Neither notify expiry nor renew automatically. Default value range: `NOTIFY_AND_MANUAL_RENEW`: Notify expiry but do not renew automatically. NOTE: it only works when charge_type is set to `PREPAID`.
* `project_id` - (Optional) ID of the project to which the instance belongs.
* `snapshot_id` - (Optional) ID of the snapshot. If specified, created the CBS by this snapshot.
* `tags` - (Optional) The available tags within this CBS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `attached` - Indicates whether the CBS is mounted the CVM.
* `storage_status` - Status of CBS. Valid values: UNATTACHED, ATTACHING, ATTACHED, DETACHING, EXPANDING, ROLLBACKING, TORECYCLE and DUMPING.


## Import

CBS storage can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_storage.storage disk-41s6jwy4
```

