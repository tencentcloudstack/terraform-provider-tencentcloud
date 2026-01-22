Use this data source to query detailed information of VPN gateways routes.

Example Usage

```hcl
data "tencentcloud_vpn_gateway_routes" "example" {
  vpn_gateway_id   = "vpngw-8dua3tbl"
  destination_cidr = "10.0.0.0/8"
  instance_type    = "VPNCONN"
  instance_id      = "vpnx-m16m4sw4"
}
```