---
subcategory: "Billing"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_billing_instance"
sidebar_current: "docs-tencentcloud-resource-billing_instance"
description: |-
  Provides a resource to create a Billing instance
---

# tencentcloud_billing_instance

Provides a resource to create a Billing instance

## Example Usage

```hcl
resource "tencentcloud_billing_instance" "example" {
  product_code     = "p_cloudfirewall"
  sub_product_code = "sp_cloudfirewall_svv1"
  region_code      = "ap-guangzhou"
  zone_code        = "ap-guangzhou-6"
  pay_mode         = "PrePay"
  parameter = jsonencode({
    "goodsNum" : 1,
    "pid" : 1002147,
    "productCode" : "p_cloudfirewall",
    "subProductCode" : "sp_cloudfirewall_svv1",
    "sv_cloudfirewall_basic_aeps" : 1,
    "sv_cloudfirewall_basic_eeps" : 0,
    "sv_cloudfirewall_basic_ipsmonth" : 0,
    "sv_cloudfirewall_basic_mon" : 0,
    "sv_cloudfirewall_basic_ueps" : 0,
    "sv_cloudfirewall_extended_ates" : 0,
    "sv_cloudfirewall_extended_clasps" : 1,
    "sv_cloudfirewall_extended_clsesps" : 0,
    "sv_cloudfirewall_extended_ex" : 0,
    "sv_cloudfirewall_extended_ibtesps" : 0,
    "sv_cloudfirewall_extended_nats" : 0,
    "sv_cloudfirewall_extended_ndr" : 0,
    "sv_cloudfirewall_extended_pcs" : 0,
    "sv_cloudfirewall_extended_spt" : 0,
    "sv_cloudfirewall_extended_sra" : 0,
    "sv_cloudfirewall_extended_srb" : 0,
    "sv_cloudfirewall_extended_sub" : 0,
    "sv_cloudfirewall_extended_subs" : 0,
    "sv_cloudfirewall_extended_vpcbges" : 0,
    "timeSpan" : 1,
    "timeUnit" : "m"
  })
  project_id  = 0
  period      = 1
  period_unit = "m"
  renew_flag  = "NOTIFY_AND_MANUAL_RENEW"
}
```

## Argument Reference

The following arguments are supported:

* `parameter` - (Required, String) Product detailed information.
* `pay_mode` - (Required, String) Payment mode. Available values: PrePay: upfront charge.
* `product_code` - (Required, String) Product code.
* `region_code` - (Required, String) Region code.
* `sub_product_code` - (Required, String) Sub-product code.
* `zone_code` - (Required, String) Availability zone code.
* `period_unit` - (Optional, String) Purchase duration unit. valid values: 
m: month,
y: year. 
default value is: m.
* `period` - (Optional, Int) Purchase duration, max number is 36, default value is 1.
* `project_id` - (Optional, Int) Project id, default value is 0.
* `renew_flag` - (Optional, String) Auto-renewal flag. valid values: NOTIFY_AND_MANUAL_RENEW: manually renew, NOTIFY_AND_AUTO_RENEW: automatically renew, DISABLE_NOTIFY_AND_MANUAL_RENEW: renewal is disabled. 
default value is NOTIFY_AND_MANUAL_RENEW.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_id` - Instance id.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.

