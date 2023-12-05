Provides a resource to create a vpc vpn_connection_reset

Example Usage

```hcl
resource "tencentcloud_vpn_connection_reset" "vpn_connection_reset" {
  vpn_gateway_id    = "vpngw-gt8bianl"
  vpn_connection_id = "vpnx-kme2tx8m"
}
```