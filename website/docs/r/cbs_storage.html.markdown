---
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
  storage_size      = "50"
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
* `storage_type` - (Required, ForceNew) Type of CBS medium, and available values include CLOUD_BASIC, CLOUD_PREMIUM and CLOUD_SSD.
* `encrypt` - (Optional, ForceNew) Indicates whether CBS is encrypted.
* `period` - (Optional) The purchased usage period of CBS, and value range [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36].
* `project_id` - (Optional) ID of the project to which the instance belongs.
* `snapshot_id` - (Optional) ID of the snapshot. If specified, created the CBS by this snapshot.
* `tags` - (Optional) The available tags within this CBS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `attached` - Indicates whether the CBS is mounted the CVM.
* `storage_status` - Status of CBS, and available values include UNATTACHED, ATTACHING, ATTACHED, DETACHING, EXPANDING, ROLLBACKING, TORECYCLE and DUMPING.


## Import

CBS storage can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_storage.storage disk-41s6jwy4
```

