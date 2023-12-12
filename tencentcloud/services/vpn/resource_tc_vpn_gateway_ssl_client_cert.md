Provides a resource to create a vpc vpn_gateway_ssl_client_cert

Example Usage

```hcl
resource "tencentcloud_vpn_gateway_ssl_client_cert" "vpn_gateway_ssl_client_cert" {
  ssl_vpn_client_id = "vpnc-123456"
  switch = "off"
}
```

Import

vpc vpn_gateway_ssl_client_cert can be imported using the id, e.g.

```
terraform import tencentcloud_vpn_gateway_ssl_client_cert.vpn_gateway_ssl_client_cert ssl_client_id
```