Provide a resource to create a VPN SSL Server.

Example Usage

```hcl
resource "tencentcloud_vpn_ssl_server" "example" {
  local_address = [
    "10.0.200.0/24",
  ]
  remote_address      = "192.168.100.0/24"
  ssl_vpn_server_name = "helloworld"
  vpn_gateway_id      = "vpngw-6lq9ayur"
#  ssl_vpn_protocol    = "UDP"
#  ssl_vpn_port        = 9798
#  integrity_algorithm = "SHA1"
#  encrypt_algorithm   = "AES-128-CBC"
#  compress            = true
}
```

Import

VPN SSL Server can be imported, e.g.

```
$ terraform import tencentcloud_vpn_ssl_server.example vpns-cik6bjct
```