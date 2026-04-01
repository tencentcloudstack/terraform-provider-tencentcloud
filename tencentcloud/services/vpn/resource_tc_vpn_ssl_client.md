Provide a resource to create a VPN SSL Client.

Example Usage

Basic Configuration

```hcl
resource "tencentcloud_vpn_ssl_client" "client" {
  ssl_vpn_server_id   = "vpns-aog5xcjj"
  ssl_vpn_client_name = "hello"
}
```

With Tags

```hcl
resource "tencentcloud_vpn_ssl_client" "client" {
  ssl_vpn_server_id   = "vpns-aog5xcjj"
  ssl_vpn_client_name = "my-ssl-client"

  # Tags for resource management
  tags = {
    Environment = "production"
    Owner       = "team-a"
  }
}
```

**Note:** Tags can be updated in-place without recreating the resource.

Import

VPN SSL Client can be imported, e.g.

```
$ terraform import tencentcloud_vpn_ssl_client.client vpn-client-id
```