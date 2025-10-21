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
* `type` - (Optional, String, ForceNew) The type of eip. Valid value:  `EIP` and `AnycastEIP` and `HighQualityEIP` and `AntiDDoSEIP`. Default is `EIP`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `public_ip` - The elastic IP address.
* `status` - The EIP current status.


## Import

EIP can be imported using the id, e.g.

```
$ terraform import tencentcloud_eip.foo eip-nyvf60va
```

