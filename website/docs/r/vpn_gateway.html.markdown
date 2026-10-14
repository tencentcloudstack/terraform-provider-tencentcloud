---
subcategory: "VPN Connections(VPN)"
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

### VPC SSL VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "example" {
  name      = "tf-example"
  bandwidth = 5
  zone      = "ap-guangzhou-3"
  type      = "SSL"
  vpc_id    = "vpc-86v957zb"

  tags = {
    createBy = "Terraform"
  }
}
```

### CCN IPSEC VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "example" {
  name      = "tf-example"
  bandwidth = 5
  zone      = "ap-guangzhou-3"
  type      = "IPSEC"

  tags = {
    createBy = "Terraform"
  }
}
```

### CCN SSL VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "example" {
  name      = "tf-example"
  bandwidth = 5
  zone      = "ap-guangzhou-3"
  type      = "SSL_CCN"

  tags = {
    createBy = "Terraform"
  }
}
```

### CCN VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "example" {
  name      = "tf-example"
  bandwidth = 200
  type      = "CCN"
  bgp_asn   = 9000

  tags = {
    createBy = "Terraform"
  }
}
```

### POSTPAID_BY_HOUR VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "example" {
  name      = "tf-example"
  vpc_id    = "vpc-dk8zmwuf"
  bandwidth = 5
  zone      = "ap-guangzhou-3"

  tags = {
    createBy = "Terraform"
  }
}
```

### PREPAID VPN gateway

```hcl
resource "tencentcloud_vpn_gateway" "example" {
  name           = "tf-example"
  vpc_id         = "vpc-dk8zmwuf"
  bandwidth      = 5
  zone           = "ap-guangzhou-3"
  charge_type    = "PREPAID"
  prepaid_period = 1

  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the VPN gateway. The length of character is limited to 1-60.
* `bandwidth` - (Optional, Int) The maximum public network output bandwidth of VPN gateway (unit: Mbps), the available values include: 5,10,20,50,100,200,500,1000. Default is 5. When charge type is `PREPAID`, bandwidth degradation operation is unsupported.
* `bgp_asn` - (Optional, Int) BGP ASN. Value range: 1 - 4294967295. Using BGP requires configuring ASN.
* `cdc_id` - (Optional, String) CDC instance ID.
* `charge_type` - (Optional, String) Charge Type of the VPN gateway. Valid value: `PREPAID`, `POSTPAID_BY_HOUR`. The default is `POSTPAID_BY_HOUR`.
* `max_connection` - (Optional, Int) Maximum number of connected clients allowed for the SSL VPN gateway. Valid values: [5, 10, 20, 50, 100]. This parameter is only required for SSL VPN gateways.
* `prepaid_period` - (Optional, Int) Period of instance to be prepaid. Valid value: `1`, `2`, `3`, `4`, `6`, `7`, `8`, `9`, `12`, `24`, `36`. The unit is month. Caution: when this para and renew_flag para are valid, the request means to renew several months more pre-paid period. This para can only be changed on `IPSEC` vpn gateway.
* `prepaid_renew_flag` - (Optional, String) Flag indicates whether to renew or not. Valid value: `NOTIFY_AND_AUTO_RENEW`, `NOTIFY_AND_MANUAL_RENEW`.
* `tags` - (Optional, Map) A list of tags used to associate different resources.
* `type` - (Optional, String) Type of gateway instance, Default is `IPSEC`. Valid value: `IPSEC`, `SSL`, `CCN` and `SSL_CCN`.
* `vpc_id` - (Optional, String, ForceNew) ID of the VPC. Required if vpn gateway is not in `CCN` or `SSL_CCN` type, and doesn't make sense for `CCN` or `SSL_CCN` vpn gateway.
* `zone` - (Optional, String, ForceNew) Zone of the VPN gateway.

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

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateVpnGateway](https://cloud.tencent.com/document/api/215/17521) - This API (CreateVpnGateway) is used to create a VPN gateway.
* Create: [DescribeVpnGateways](https://cloud.tencent.com/document/api/215/17514) - This API (DescribeVpnGateways) is used to query the VPN gateway list.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Read: [DescribeVpnGateways](https://cloud.tencent.com/document/api/215/17514) - This API (DescribeVpnGateways) is used to query the VPN gateway list.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Update: [ModifyVpnGatewayAttribute](https://cloud.tencent.com/document/api/215/17507) - This API (ModifyVpnGatewayAttribute) is used to modify the attributes of VPN gateways.
* Update: [RenewVpnGateway](https://cloud.tencent.com/document/api/215/17506) - This API (RenewVpnGateway) is used to renew prepaid (yearly/monthly subscription) VPN gateways. Currently, only IPSEC gateways are supported.
* Update: [ResetVpnGatewayInternetMaxBandwidth](https://cloud.tencent.com/document/api/215/17504) - This API is used to adjust the bandwidth cap of a VPN gateway. The adjustment of the VPN gateway bandwidth is limited to [5,100] Mbps and [200,1000] Mbps.
* Update: [SetVpnGatewaysRenewFlag](https://cloud.tencent.com/document/api/215/81120) - This API is used set the auto-renewal configuration of a VPN gateway.
* Delete: [DeleteVpnGateway](https://cloud.tencent.com/document/api/215/17518) - This API (DeleteVpnGateway) is used to delete a VPN gateway. Currently, only deletion of pay-as-you-go IPSEC gateway instances in running status is supported.
* Delete: [DescribeVpnConnections](https://cloud.tencent.com/document/api/215/17515) - This API is used to used to query the list of VPN tunnels.
* Delete: [DescribeVpnGateways](https://cloud.tencent.com/document/api/215/17514) - This API (DescribeVpnGateways) is used to query the VPN gateway list.

## Import

VPN gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_gateway.example vpngw-8ccsnclt
```

