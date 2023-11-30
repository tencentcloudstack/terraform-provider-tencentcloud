Provides a resource to create a vpn_gateway_ccn_routes

Example Usage

```hcl
resource "tencentcloud_vpn_gateway_ccn_routes" "vpn_gateway_ccn_routes" {
  destination_cidr_block = "192.168.1.0/24"
  route_id               = "vpnr-akdy0757"
  status                 = "DISABLE"
  vpn_gateway_id         = "vpngw-lie1a4u7"
}

```

Import

vpc vpn_gateway_ccn_routes can be imported using the id, e.g.

```
terraform import tencentcloud_vpn_gateway_ccn_routes.vpn_gateway_ccn_routes vpn_gateway_id#ccn_routes_id
```