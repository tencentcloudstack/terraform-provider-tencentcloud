---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip"
sidebar_current: "docs-tencentcloud-resource-eip"
description: |-
  Provides an EIP resource.
---

# tencentcloud_eip

Provides an EIP resource.

## Example Usage

### Paid by the bandwidth package

```hcl
resource "tencentcloud_eip" "foo" {
  name                 = "awesome_gateway_ip"
  bandwidth_package_id = "bwp-jtvzuky6"
  internet_charge_type = "BANDWIDTH_PACKAGE"
  type                 = "EIP"
}
```

### AntiDDos Eip

```hcl
resource "tencentcloud_eip" "foo" {
  name                 = "awesome_gateway_ip"
  bandwidth_package_id = "bwp-4ocyia9s"
  internet_charge_type = "BANDWIDTH_PACKAGE"
  type                 = "AntiDDoSEIP"
  anti_ddos_package_id = "xxxxxxxx"

  tags = {
    "test" = "test"
  }
}
```

### Eip With Network Egress

```hcl
resource "tencentcloud_eip" "foo" {
  name                       = "egress_eip"
  egress                     = "center_egress2"
  internet_charge_type       = "BANDWIDTH_PACKAGE"
  internet_service_provider  = "CMCC"
  internet_max_bandwidth_out = 1
  type                       = "EIP"
}
```

## Argument Reference

The following arguments are supported:

* `anti_ddos_package_id` - (Optional, String) ID of anti DDos package, it must set when `type` is `AntiDDoSEIP`.
* `anycast_zone` - (Optional, String, ForceNew) The zone of anycast. Valid value: `ANYCAST_ZONE_GLOBAL` and `ANYCAST_ZONE_OVERSEAS`.
* `applicable_for_clb` - (Optional, Bool, **Deprecated**) It has been deprecated from version 1.27.0. Indicates whether the anycast eip can be associated to a CLB.
* `auto_renew_flag` - (Optional, Int) Auto renew flag.  0 - default state (manual renew); 1 - automatic renew; 2 - explicit no automatic renew. NOTES: Only supported prepaid EIP.
* `bandwidth_package_id` - (Optional, String) ID of bandwidth package, it will set when `internet_charge_type` is `BANDWIDTH_PACKAGE`.
* `cdc_id` - (Optional, String) CDC Unique ID.
* `egress` - (Optional, String) Network egress. It defaults to `center_egress1`. If you want to try the egress feature, please [submit a ticket](https://console.cloud.tencent.com/workorder/category).
* `internet_charge_type` - (Optional, String) The charge type of eip. Valid values: `BANDWIDTH_PACKAGE`, `BANDWIDTH_POSTPAID_BY_HOUR`, `BANDWIDTH_PREPAID_BY_MONTH` and `TRAFFIC_POSTPAID_BY_HOUR`.
* `internet_max_bandwidth_out` - (Optional, Int) The bandwidth limit of EIP, unit is Mbps.
* `internet_service_provider` - (Optional, String, ForceNew) Internet service provider of eip. Valid value: `BGP`, `CMCC`, `CTCC` and `CUCC`.
* `name` - (Optional, String) The name of eip.
* `prepaid_period` - (Optional, Int) Period of instance. Default value: `1`. Valid value: `1`, `2`, `3`, `4`, `6`, `7`, `8`, `9`, `12`, `24`, `36`. NOTES: must set when `internet_charge_type` is `BANDWIDTH_PREPAID_BY_MONTH`.
* `tags` - (Optional, Map) The tags of eip.
* `type` - (Optional, String, ForceNew) The type of eip. Valid value:  `EIP` and `AnycastEIP` and `HighQualityEIP` and `AntiDDoSEIP` and `ResidentialEIP`. Default is `EIP`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `public_ip` - The elastic IP address.
* `status` - The EIP current status.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AllocateAddresses](https://cloud.tencent.com/document/api/215/16699) - This API is used to apply for one or more [Elastic IP Addresses](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1) (EIPs for short).
* Create: [DescribeAddresses](https://cloud.tencent.com/document/api/215/16702) - This API (DescribeAddresses) is used to query the information of one or multiple [Elastic IPs](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1).
* Create: [DescribeTaskResult](https://cloud.tencent.com/document/api/215/36271) - This API is used to query the EIP async job execution results.
* Read: [DescribeAddresses](https://cloud.tencent.com/document/api/215/16702) - This API (DescribeAddresses) is used to query the information of one or multiple [Elastic IPs](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1).
* Read: [DescribeBandwidthPackages](https://cloud.tencent.com/document/api/215/19209) - This API is used to query bandwidth package information, including the unique ID of the bandwidth package, the type, the billing mode, the name, and the resource information.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Update: [ModifyAddressAttribute](https://cloud.tencent.com/document/api/215/16704) - This API (ModifyAddressAttribute) is used to modify the name of an [Elastic IP](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1).
* Update: [ModifyAddressInternetChargeType](https://cloud.tencent.com/document/api/215/43849) - This API is used to adjust the network billing mode of an EIP. Please note that it's available to users whose network fees are billed on IPs but not CVMs.
* Update: [ModifyAddressesBandwidth](https://cloud.tencent.com/document/api/215/19214) - This API is used to adjust the bandwidth of [Elastic IP](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1), including EIP billed on a pay-as-you-go, monthly subscription, and bandwidth package basis.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Update: [RenewAddresses](https://cloud.tencent.com/document/api/215/49015)
* Delete: [DescribeAddresses](https://cloud.tencent.com/document/api/215/16702) - This API (DescribeAddresses) is used to query the information of one or multiple [Elastic IPs](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1).
* Delete: [DescribeNatGateways](https://cloud.tencent.com/document/api/215/36034) - This API is used to query NAT gateways.
* Delete: [DescribeTaskResult](https://cloud.tencent.com/document/api/215/36271) - This API is used to query the EIP async job execution results.
* Delete: [DisassociateAddress](https://cloud.tencent.com/document/api/215/16703) - This API is used to unbind an [Elastic IP](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1) (EIP for short).
* Delete: [DisassociateNatGatewayAddress](https://cloud.tencent.com/document/api/215/36716) - This API is used to unbind an EIP from a NAT gateway.
* Delete: [ReleaseAddresses](https://cloud.tencent.com/document/api/215/16705) - This API (ReleaseAddresses) is used to release one or multiple [Elastic IPs](https://intl.cloud.tencent.com/document/product/213/1941?from_cn_redirect=1).

## Import

EIP can be imported using the id, e.g.

```
$ terraform import tencentcloud_eip.foo eip-nyvf60va
```

