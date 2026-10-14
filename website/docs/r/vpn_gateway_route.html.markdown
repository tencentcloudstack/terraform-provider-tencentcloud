---
subcategory: "VPN Connections(VPN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_gateway_route"
sidebar_current: "docs-tencentcloud-resource-vpn_gateway_route"
description: |-
  Provides a resource to create a VPN gateway route.
---

# tencentcloud_vpn_gateway_route

Provides a resource to create a VPN gateway route.

## Example Usage

```hcl
resource "tencentcloud_vpn_gateway_route" "route" {
  vpn_gateway_id         = "vpngw-ak9sjem2"
  destination_cidr_block = "10.0.0.0/16"
  instance_id            = "vpnx-5b5dmao3"
  instance_type          = "VPNCONN"
  priority               = 100
  status                 = "DISABLE"
}
```

## Argument Reference

The following arguments are supported:

* `destination_cidr_block` - (Required, String, ForceNew) Destination IDC IP range.
* `instance_id` - (Required, String, ForceNew) Instance ID of the next hop.
* `instance_type` - (Required, String, ForceNew) Next hop type (type of the associated instance). Valid values: VPNCONN (VPN tunnel) and CCN (CCN instance).
* `priority` - (Required, Int, ForceNew) Priority. Valid values: 0 and 100.
* `status` - (Required, String) Status. Valid values: ENABLE and DISABLE.
* `vpn_gateway_id` - (Required, String, ForceNew) VPN gateway ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `route_id` - Route ID.
* `type` - Route type. Default value: Static.
* `update_time` - Update time.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateVpnGatewayRoutes](https://cloud.tencent.com/document/api/215/57678) - This API is used to create destination routes of a route-based VPN gateway.
* Read: [DescribeVpnGatewayRoutes](https://cloud.tencent.com/document/api/215/57676) - This API is used to query VPN gateway routes.
* Update: [ModifyVpnGatewayRoutes](https://cloud.tencent.com/document/api/215/57675) - This API is used to modify VPN gateway routes.
* Delete: [DeleteVpnGatewayRoutes](https://cloud.tencent.com/document/api/215/57677) - This API is used to delete routes of a VPN gateway.

## Import

VPN gateway route can be imported using the id, the id format must be '{vpn_gateway_id}#{route_id}', e.g.

```
$ terraform import tencentcloud_vpn_gateway_route.route1 vpngw-ak9sjem2#vpngw-8ccsnclt
```

