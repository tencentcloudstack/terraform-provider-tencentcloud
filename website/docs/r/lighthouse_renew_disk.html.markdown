---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_renew_disk"
sidebar_current: "docs-tencentcloud-resource-lighthouse_renew_disk"
description: |-
  Provides a resource to create a lighthouse renew_disk
---

# tencentcloud_lighthouse_renew_disk

Provides a resource to create a lighthouse renew_disk

## Example Usage

```hcl
resource "tencentcloud_lighthouse_renew_disk" "renew_disk" {
  disk_id = "lhdisk-xxxxxx"
  renew_disk_charge_prepaid {
    period     = 1
    renew_flag = "NOTIFY_AND_AUTO_RENEW"
    time_unit  = "m"
  }
  auto_voucher = true
}
```

## Argument Reference

The following arguments are supported:

* `disk_id` - (Required, String, ForceNew) List of disk ID.
* `renew_disk_charge_prepaid` - (Required, List, ForceNew) Renew cloud hard disk subscription related parameter settings.
* `auto_voucher` - (Optional, Bool, ForceNew) Whether to automatically use the voucher. Not used by default.

The `renew_disk_charge_prepaid` object supports the following:

* `cur_instance_deadline` - (Optional, String) Current instance expiration time. Such as 2018-01-01 00:00:00. Specifying this parameter can align the expiration time of the instance attached to the disk. One of this parameter and Period must be specified, and cannot be specified at the same time.
* `period` - (Optional, Int) Renewal period.
* `renew_flag` - (Optional, String) Automatic renewal falg. Value:NOTIFY_AND_AUTO_RENEW: Notice expires and auto-renews.NOTIFY_AND_MANUAL_RENEW: Notification expires without automatic renewal, users need to manually renew.DISABLE_NOTIFY_AND_AUTO_RENEW: No automatic renewal and no notification.Default: NOTIFY_AND_MANUAL_RENEW. If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the disk will be automatically renewed monthly when the account balance is sufficient.
* `time_unit` - (Optional, String) newly purchased unit. Default: m.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



