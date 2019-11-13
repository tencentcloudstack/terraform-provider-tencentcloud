---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_gateway"
sidebar_current: "docs-tencentcloud-resource-vpn_gateway"
description: |-
  Provides a resource to create a VPN gateway.
---

# tencentcloud_vpn_gateway

Provides a resource to create a VPN gateway.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the VPN gateway. The length of character is limited to 1-60.
* `vpc_id` - (Required, ForceNew) ID of the VPC.
* `zone` - (Required, ForceNew) Zone of the VPN gateway.
* `bandwidth` - (Optional) The maximum public network output bandwidth of VPN gateway (unit: Mbps), the available values include: 5,10,20,50,100. Default is 5.
* `tags` - (Optional) A list of tags used to associate different resources.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `charge_type` - Charge Type of the VPN gateway, valid values are `PREPAID`, `POSTPAID_BY_HOUR` and default is `POSTPAID_BY_HOUR`.
* `create_time` - Create time of the VPN gateway.
* `expired_time` - Expired time of the VPN gateway when charge type is `PREPAID`.
* `is_address_blocked` - Indicates whether ip address is blocked.
* `new_purchase_plan` - The plan of new purchase, valid value is `PREPAID_TO_POSTPAID`.
* `prepaid_period` - Period of instance to be prepaid. Valid values are 1, 2, 3, 4, 6, 7, 8, 9, 12, 24, 36 and unit is month.
* `prepaid_renew_flag` - Flag indicates whether to renew or not, valid values are `NOTIFY_AND_RENEW`, `NOTIFY_AND_AUTO_RENEW`, `NOT_NOTIFY_AND_NOT_RENEW`.
* `public_ip_address` - Public ip of the VPN gateway.
* `restrict_state` - Restrict state of gateway, valid values are `PRETECIVELY_ISOLATED`, `NORMAL`.
* `state` - State of the VPN gateway, valid values are `PENDING`, `DELETING`, `AVAILABLE`.
* `type` - Type of gateway instance, valid values are `IPSEC`, `SSL`.


## Import

VPN gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_gateway.foo vpngw-8ccsnclt
```

