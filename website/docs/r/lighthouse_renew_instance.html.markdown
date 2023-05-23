---
subcategory: "TencentCloud Lighthouse(Lighthouse)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_lighthouse_renew_instance"
sidebar_current: "docs-tencentcloud-resource-lighthouse_renew_instance"
description: |-
  Provides a resource to create a lighthouse renew_instance
---

# tencentcloud_lighthouse_renew_instance

Provides a resource to create a lighthouse renew_instance

## Example Usage

```hcl
resource "tencentcloud_lighthouse_renew_instance" "renew_instance" {
  instance_id =
  instance_charge_prepaid {
    period     = 1
    renew_flag = "NOTIFY_AND_MANUAL_RENEW"

  }
  renew_data_disk = true
  auto_voucher    = false
}
```

## Argument Reference

The following arguments are supported:

* `instance_charge_prepaid` - (Required, List, ForceNew) Prepaid mode, that is, yearly and monthly subscription related parameter settings. Through this parameter, you can specify attributes such as the purchase duration of the Subscription instance and whether to set automatic renewal.
* `instance_id` - (Required, String, ForceNew) Instance ID.
* `auto_voucher` - (Optional, Bool, ForceNew) Whether to automatically deduct vouchers. Valid values:
- true: Automatically deduct vouchers.
-false:Do not automatically deduct vouchers. Default value: false.
* `renew_data_disk` - (Optional, Bool, ForceNew) Whether to renew the data disk. Valid values:true: Indicates that the renewal instance also renews the data disk attached to it.false: Indicates that the instance will be renewed and the data disk attached to it will not be renewed at the same time.Default value: true.

The `instance_charge_prepaid` object supports the following:

* `period` - (Required, Int) The duration of purchasing an instance. Unit is month. Valid values are (1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60).
* `renew_flag` - (Optional, String) Automatic renewal logo. Values:
- `NOTIFY_AND_AUTO_RENEW`: notify expiration and renew automatically;
- `NOTIFY_AND_MANUAL_RENEW`: notification of expiration does not renew automatically. Users need to renew manually;
- `DISABLE_NOTIFY_AND_AUTO_RENEW`: no automatic renewal and no notification;
Default value: `NOTIFY_AND_MANUAL_RENEW`. If this parameter is specified as `NOTIFY_AND_AUTO_RENEW`, the instance will be automatically renewed on a monthly basis after expiration, when the account balance is sufficient.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



