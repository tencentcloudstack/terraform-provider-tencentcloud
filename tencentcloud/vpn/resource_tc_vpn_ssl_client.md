Provide a resource to create a VPN SSL Client.

Example Usage

```hcl
resource "tencentcloud_vpn_ssl_client" "client" {
  ssl_vpn_server_id = "vpns-aog5xcjj"
  ssl_vpn_client_name = "hello"
}

```

Import

VPN SSL Client can be imported, e.g.

```
$ terraform import tencentcloud_vpn_ssl_client.client vpn-client-id
```