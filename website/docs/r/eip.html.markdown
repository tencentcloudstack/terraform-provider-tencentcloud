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

```hcl
resource "tencentcloud_eip" "foo" {
  name = "awesome_gateway_ip"
}
```

## Argument Reference

The following arguments are supported:

* `anycast_zone` - (Optional, String, ForceNew) The zone of anycast. Valid value: `ANYCAST_ZONE_GLOBAL` and `ANYCAST_ZONE_OVERSEAS`.
* `applicable_for_clb` - (Optional, Bool, **Deprecated**) It has been deprecated from version 1.27.0. Indicates whether the anycast eip can be associated to a CLB.
* `internet_charge_type` - (Optional, String, ForceNew) The charge type of eip. Valid values: `BANDWIDTH_PACKAGE`, `BANDWIDTH_POSTPAID_BY_HOUR`, `BANDWIDTH_PREPAID_BY_MONTH` and `TRAFFIC_POSTPAID_BY_HOUR`.
* `internet_max_bandwidth_out` - (Optional, Int) The bandwidth limit of EIP, unit is Mbps.
* `internet_service_provider` - (Optional, String, ForceNew) Internet service provider of eip. Valid value: `BGP`, `CMCC`, `CTCC` and `CUCC`.
* `name` - (Optional, String) The name of eip.
* `tags` - (Optional, Map) The tags of eip.
* `type` - (Optional, String, ForceNew) The type of eip. Valid value:  `EIP` and `AnycastEIP` and `HighQualityEIP`. Default is `EIP`.

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

