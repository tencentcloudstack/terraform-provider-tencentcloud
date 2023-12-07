Provides a resource to create a vpc vpn_customer_gateway_configuration_download

Example Usage

```hcl
resource "tencentcloud_vpn_customer_gateway_configuration_download" "vpn_customer_gateway_configuration_download" {
  vpn_gateway_id    = "vpngw-gt8bianl"
  vpn_connection_id = "vpnx-kme2tx8m"
  customer_gateway_vendor {
    platform         = "comware"
    software_version = "V1.0"
    vendor_name      = "h3c"
  }
  interface_name    = "test"
}
```