---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_renew_instance"
sidebar_current: "docs-tencentcloud-resource-cvm_renew_instance"
description: |-
  Provides a resource to create a cvm renew_instance
---

# tencentcloud_cvm_renew_instance

Provides a resource to create a cvm renew_instance

## Example Usage

```hcl
resource "tencentcloud_cvm_renew_instance" "renew_instance" {
  instance_ids =
  instance_charge_prepaid {
    period     = 1
    renew_flag = "NOTIFY_AND_AUTO_RENEW"
  }
  renew_portable_data_disk = true
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance ID.
* `instance_charge_prepaid` - (Optional, List, ForceNew) Prepaid mode, that is, yearly and monthly subscription related parameter settings. Through this parameter, you can specify the renewal duration of the Subscription instance, whether to set automatic renewal, and other attributes. For yearly and monthly subscription instances, this parameter is required.
* `renew_portable_data_disk` - (Optional, Bool, ForceNew) Whether to renew the elastic data disk. Valid values:
- `TRUE`: Indicates to renew the subscription instance and renew the attached elastic data disk at the same time
- `FALSE`: Indicates that the subscription instance will be renewed and the elastic data disk attached to it will not be renewed
Default value: TRUE.

The `instance_charge_prepaid` object supports the following:

* `period` - (Required, Int) Subscription period; unit: month; valid values: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36, 48, 60. Note: This field may return null, indicating that no valid value is found.
* `renew_flag` - (Optional, String) Auto renewal flag. Valid values:
- `NOTIFY_AND_AUTO_RENEW`: notify upon expiration and renew automatically;
- `NOTIFY_AND_MANUAL_RENEW`: notify upon expiration but do not renew automatically;
- `DISABLE_NOTIFY_AND_MANUAL_RENEW`: neither notify upon expiration nor renew automatically;
Default value: NOTIFY_AND_MANUAL_RENEW. If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis if the account balance is sufficient. Note: This field may return null, indicating that no valid value is found.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



