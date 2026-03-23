Provide a resource to create a VPN SSL Server.

Example Usage

Basic Configuration

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

With Tags and DNS Configuration

```hcl
resource "tencentcloud_vpn_ssl_server" "example" {
  local_address = [
    "10.0.200.0/24",
  ]
  remote_address      = "192.168.100.0/24"
  ssl_vpn_server_name = "helloworld"
  vpn_gateway_id      = "vpngw-6lq9ayur"

  # Tags for resource management
  tags = {
    Environment = "production"
    Owner       = "team-a"
  }

  # Custom DNS servers
  dns_servers {
    primary_dns   = "8.8.8.8"
    secondary_dns = "8.8.4.4"
  }
}
```

With SSO Authentication (Requires Whitelist)

**Note:** SSO authentication feature requires whitelist approval from TencentCloud. Please contact TencentCloud support to apply for whitelist access before enabling this feature.

```hcl
resource "tencentcloud_vpn_ssl_server" "example" {
  local_address = [
    "10.0.200.0/24",
  ]
  remote_address      = "192.168.100.0/24"
  ssl_vpn_server_name = "helloworld"
  vpn_gateway_id      = "vpngw-6lq9ayur"

  # Enable SSO authentication
  sso_enabled = true
  saml_data   = "<SAML configuration data>" # Replace with your SAML data
}
```

With Access Policy Control

**Note:** The `access_policy_enabled` parameter only controls the feature switch. Detailed access policies must be configured through the TencentCloud console or other resources.

```hcl
resource "tencentcloud_vpn_ssl_server" "example" {
  local_address = [
    "10.0.200.0/24",
  ]
  remote_address      = "192.168.100.0/24"
  ssl_vpn_server_name = "helloworld"
  vpn_gateway_id      = "vpngw-6lq9ayur"

  # Enable access policy control
  access_policy_enabled = true
}
```

Import

VPN SSL Server can be imported, e.g.

```
$ terraform import tencentcloud_vpn_ssl_server.example vpns-cik6bjct
```