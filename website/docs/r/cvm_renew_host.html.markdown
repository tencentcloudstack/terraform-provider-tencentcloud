---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_renew_host"
sidebar_current: "docs-tencentcloud-resource-cvm_renew_host"
description: |-
  Provides a resource to create a cvm renew_host
---

# tencentcloud_cvm_renew_host

Provides a resource to create a cvm renew_host

## Example Usage

```hcl
resource "tencentcloud_cvm_renew_host" "renew_host" {
  host_id = "xxxxxx"
  host_charge_prepaid {
    period     = 1
    renew_flag = "NOTIFY_AND_MANUAL_RENEW"
  }
}
```

## Argument Reference

The following arguments are supported:

* `host_charge_prepaid` - (Required, List, ForceNew) Prepaid mode, that is, yearly and monthly subscription related parameter settings. Through this parameter, you can specify attributes such as the purchase duration of the Subscription instance and whether to set automatic renewal. If the payment mode of the specified instance is prepaid, this parameter must be passed.
* `host_id` - (Required, String, ForceNew) CDH instance ID.

The `host_charge_prepaid` object supports the following:

* `period` - (Required, Int) The duration of purchasing an instance, unit: month. Value range: 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 24, 36.
* `renew_flag` - (Optional, String) Auto renewal flag. Valid values:&lt;br&gt;&lt;li&gt;NOTIFY_AND_AUTO_RENEW: notify upon expiration and renew automatically&lt;br&gt;&lt;li&gt;NOTIFY_AND_MANUAL_RENEW: notify upon expiration but do not renew automatically&lt;br&gt;&lt;li&gt;DISABLE_NOTIFY_AND_MANUAL_RENEW: neither notify upon expiration nor renew automatically&lt;br&gt;&lt;br&gt;Default value: NOTIFY_AND_AUTO_RENEW. If this parameter is specified as NOTIFY_AND_AUTO_RENEW, the instance will be automatically renewed on a monthly basis if the account balance is sufficient.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [RenewHosts](https://www.tencentcloud.com/search/?wd=RenewHosts&type=document&lang=en)

