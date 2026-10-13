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

* Create: [AllocateAddresses](https://cloud.tencent.com/document/api/215/16699) - 创建弹性公网IP
* Create: [DescribeAddresses](https://cloud.tencent.com/document/api/215/16702) - 查询弹性公网IP列表
* Create: [DescribeTaskResult](https://cloud.tencent.com/document/api/215/36271) - 查询异步任务执行结果
* Read: [DescribeAddresses](https://cloud.tencent.com/document/api/215/16702) - 查询弹性公网IP列表
* Read: [DescribeBandwidthPackages](https://cloud.tencent.com/document/api/215/19209) - 查询带宽包资源
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - 批量查看资源关联的标签
* Update: [ModifyAddressAttribute](https://cloud.tencent.com/document/api/215/16704) - 修改弹性公网IP属性
* Update: [ModifyAddressInternetChargeType](https://cloud.tencent.com/document/api/215/43849) - 调整弹性公网ip计费模式
* Update: [ModifyAddressesBandwidth](https://cloud.tencent.com/document/api/215/19214) - 调整弹性公网IP带宽
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - 批量修改资源关联的标签
* Update: [RenewAddresses](https://cloud.tencent.com/document/api/215/49015) - 续费包月带宽弹性公网IP
* Delete: [DescribeAddresses](https://cloud.tencent.com/document/api/215/16702) - 查询弹性公网IP列表
* Delete: [DescribeNatGateways](https://cloud.tencent.com/document/api/215/36034) - 查询NAT网关
* Delete: [DescribeTaskResult](https://cloud.tencent.com/document/api/215/36271) - 查询异步任务执行结果
* Delete: [DisassociateAddress](https://cloud.tencent.com/document/api/215/16703) - 解绑弹性公网IP
* Delete: [DisassociateNatGatewayAddress](https://cloud.tencent.com/document/api/215/36716) - NAT网关解绑弹性IP
* Delete: [ReleaseAddresses](https://cloud.tencent.com/document/api/215/16705) - 释放弹性公网IP

## Import

EIP can be imported using the id, e.g.

```
$ terraform import tencentcloud_eip.foo eip-nyvf60va
```

