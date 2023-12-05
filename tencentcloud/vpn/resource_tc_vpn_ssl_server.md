Provide a resource to create a VPN SSL Server.

Example Usage

```hcl
resource "tencentcloud_vpn_ssl_server" "server" {
  local_address       = [
    "10.0.0.0/17",
  ]
  remote_address      = "11.0.0.0/16"
  ssl_vpn_server_name = "helloworld"
  vpn_gateway_id      = "vpngw-335lwf7d"
  ssl_vpn_protocol = "UDP"
  ssl_vpn_port = 1194
  integrity_algorithm = "MD5"
  encrypt_algorithm = "AES-128-CBC"
  compress = true
}
```

Import

VPN SSL Server can be imported, e.g.

```
$ terraform import tencentcloud_vpn_ssl_server.server vpn-server-id
```