---
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

* `anycast_zone` - (Optional, ForceNew) The zone of anycast, and available values include `ANYCAST_ZONE_GLOBAL` and `ANYCAST_ZONE_OVERSEAS`.
* `applicable_for_clb` - (Optional, ForceNew, **Deprecated**) It has been deprecated from version 1.27.0. Indicates whether the anycast eip can be associated to a CLB.
* `internet_charge_type` - (Optional, ForceNew) The charge type of eip, and available values include `BANDWIDTH_PACKAGE`, `BANDWIDTH_POSTPAID_BY_HOUR` and `TRAFFIC_POSTPAID_BY_HOUR`.
* `internet_max_bandwidth_out` - (Optional, ForceNew) The bandwidth limit of eip, unit is Mbps, and the range is 1-1000.
* `internet_service_provider` - (Optional, ForceNew) Internet service provider of eip, and available values include `BGP`, `CMCC`, `CTCC` and `CUCC`.
* `name` - (Optional) The name of eip.
* `tags` - (Optional) The tags of eip.
* `type` - (Optional, ForceNew) The type of eip, and available values include `EIP` and `AnycastEIP`. Default is `EIP`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `public_ip` - The elastic ip address.
* `status` - The eip current status.


## Import

EIP can be imported using the id, e.g.

```
$ terraform import tencentcloud_eip.foo eip-nyvf60va
```

