---
subcategory: "VPN"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_gateway"
sidebar_current: "docs-tencentcloud-resource-vpn_gateway"
description: |-
  Provides a resource to create a VPN gateway.
---

# tencentcloud_vpn_gateway

Provides a resource to create a VPN gateway.

-> **NOTE:** The prepaid VPN gateway do not support renew operation or delete operation with terraform.

## Example Usage

POSTPAID_BY_HOUR VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "my_cgw" {
  name      = "test"
  vpc_id    = "vpc-dk8zmwuf"
  bandwidth = 5
  zone      = "ap-guangzhou-3"

  tags = {
    test = "test"
  }
}
```

PREPAID VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "my_cgw" {
  name           = "test"
  vpc_id         = "vpc-dk8zmwuf"
  bandwidth      = 5
  zone           = "ap-guangzhou-3"
  charge_type    = "PREPAID"
  prepaid_period = 1

  tags = {
    test = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the VPN gateway. The length of character is limited to 1-60.
* `zone` - (Required, ForceNew) Zone of the VPN gateway.
* `bandwidth` - (Optional) The maximum public network output bandwidth of VPN gateway (unit: Mbps), the available values include: 5,10,20,50,100,200,500,1000. Default is 5. When charge type is `PREPAID`, bandwidth degradation operation is unsupported.
* `cdc_id` - (Optional) CDC instance ID.
* `charge_type` - (Optional) Charge Type of the VPN gateway. Valid value: `PREPAID`, `POSTPAID_BY_HOUR`. The default is `POSTPAID_BY_HOUR`.
* `max_connection` - (Optional) Maximum number of connected clients allowed for the SSL VPN gateway. Valid values: [5, 10, 20, 50, 100]. This parameter is only required for SSL VPN gateways.
* `prepaid_period` - (Optional) Period of instance to be prepaid. Valid value: `1`, `2`, `3`, `4`, `6`, `7`, `8`, `9`, `12`, `24`, `36`. The unit is month. Caution: when this para and renew_flag para are valid, the request means to renew several months more pre-paid period. This para can only be set to take effect in create operation.
* `prepaid_renew_flag` - (Optional) Flag indicates whether to renew or not. Valid value: `NOTIFY_AND_RENEW`, `NOTIFY_AND_AUTO_RENEW`, `NOT_NOTIFY_AND_NOT_RENEW`. This para can only be set to take effect in create operation.
* `tags` - (Optional) A list of tags used to associate different resources.
* `type` - (Optional) Type of gateway instance. Valid value: `IPSEC`, `SSL` and `CCN`. Note: CCN type is only for whitelist customer now.
* `vpc_id` - (Optional, ForceNew) ID of the VPC. Required if vpn gateway is not in `CCN` type, and doesn't make sense for `CCN` vpn gateway.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the VPN gateway.
* `expired_time` - Expired time of the VPN gateway when charge type is `PREPAID`.
* `is_address_blocked` - Indicates whether ip address is blocked.
* `new_purchase_plan` - The plan of new purchase. Valid value: `PREPAID_TO_POSTPAID`.
* `public_ip_address` - Public IP of the VPN gateway.
* `restrict_state` - Restrict state of gateway. Valid value: `PRETECIVELY_ISOLATED`, `NORMAL`.
* `state` - State of the VPN gateway. Valid value: `PENDING`, `DELETING`, `AVAILABLE`.


## Import

VPN gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_gateway.foo vpngw-8ccsnclt
```

