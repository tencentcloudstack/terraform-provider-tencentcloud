---
subcategory: "VPN Connections(VPN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_gateway_ccn_routes"
sidebar_current: "docs-tencentcloud-resource-vpn_gateway_ccn_routes"
description: |-
  Provides a resource to create a vpn_gateway_ccn_routes
---

# tencentcloud_vpn_gateway_ccn_routes

Provides a resource to create a vpn_gateway_ccn_routes

## Example Usage

```hcl
resource "tencentcloud_vpn_gateway_ccn_routes" "vpn_gateway_ccn_routes" {
  destination_cidr_block = "192.168.1.0/24"
  route_id               = "vpnr-akdy0757"
  status                 = "DISABLE"
  vpn_gateway_id         = "vpngw-lie1a4u7"
}
```

## Argument Reference

The following arguments are supported:

* `destination_cidr_block` - (Required, String, ForceNew) Routing CIDR.
* `route_id` - (Required, String, ForceNew) Route Id.
* `status` - (Required, String) Whether routing information is enabled. `ENABLE`: Enable Route, `DISABLE`: Disable Route.
* `vpn_gateway_id` - (Required, String, ForceNew) VPN GATEWAY INSTANCE ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc vpn_gateway_ccn_routes can be imported using the id, e.g.

```
terraform import tencentcloud_vpn_gateway_ccn_routes.vpn_gateway_ccn_routes vpn_gateway_id#ccn_routes_id
```

