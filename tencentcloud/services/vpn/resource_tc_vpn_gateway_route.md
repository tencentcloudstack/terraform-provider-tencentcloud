Provides a resource to create a VPN gateway route.

Example Usage

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

Import

VPN gateway route can be imported using the id, the id format must be '{vpn_gateway_id}#{route_id}', e.g.

```
$ terraform import tencentcloud_vpn_gateway_route.route1 vpngw-ak9sjem2#vpngw-8ccsnclt
```