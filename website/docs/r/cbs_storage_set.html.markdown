---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_storage_set"
sidebar_current: "docs-tencentcloud-resource-cbs_storage_set"
description: |-
  Provides a resource to create CBS set.
---

# tencentcloud_cbs_storage_set

Provides a resource to create CBS set.

## Example Usage

```hcl
resource "tencentcloud_cbs_storage_set" "storage" {
  disk_count        = 10
  storage_name      = "mystorage"
  storage_type      = "CLOUD_SSD"
  storage_size      = 100
  availability_zone = "ap-guangzhou-3"
  project_id        = 0
  encrypt           = false
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, String, ForceNew) The available zone that the CBS instance locates at.
* `storage_name` - (Required, String) Name of CBS. The maximum length can not exceed 60 bytes.
* `storage_size` - (Required, Int) Volume of CBS, and unit is GB.
* `storage_type` - (Required, String, ForceNew) Type of CBS medium. Valid values: CLOUD_BASIC: HDD cloud disk, CLOUD_PREMIUM: Premium Cloud Storage, CLOUD_BSSD: General Purpose SSD, CLOUD_SSD: SSD, CLOUD_HSSD: Enhanced SSD, CLOUD_TSSD: Tremendous SSD.
* `charge_type` - (Optional, String) The charge type of CBS instance. Only support `POSTPAID_BY_HOUR`.
* `disk_count` - (Optional, Int, ForceNew) The number of disks to be purchased. Default 1.
* `encrypt` - (Optional, Bool, ForceNew) Indicates whether CBS is encrypted.
* `project_id` - (Optional, Int) ID of the project to which the instance belongs.
* `snapshot_id` - (Optional, String) ID of the snapshot. If specified, created the CBS by this snapshot.
* `throughput_performance` - (Optional, Int) Add extra performance to the data disk. Only works when disk type is `CLOUD_TSSD` or `CLOUD_HSSD`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `attached` - Indicates whether the CBS is mounted the CVM.
* `disk_ids` - disk id list.
* `storage_status` - Status of CBS. Valid values: UNATTACHED, ATTACHING, ATTACHED, DETACHING, EXPANDING, ROLLBACKING, TORECYCLE and DUMPING.


