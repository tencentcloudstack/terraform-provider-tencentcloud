Use this data source to query detailed information of VPN connections.

Example Usage

```hcl
data "tencentcloud_vpn_connections" "foo" {
  name                = "main"
  id                  = "vpnx-xfqag"
  vpn_gateway_id      = "vpngw-8ccsnclt"
  vpc_id              = "cgw-xfqag"
  customer_gateway_id = ""
  tags = {
    test = "tf"
  }
}
```