---
subcategory: "VPN Connections(VPN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_ssl_server"
sidebar_current: "docs-tencentcloud-resource-vpn_ssl_server"
description: |-
  Provide a resource to create a VPN SSL Server.
---

# tencentcloud_vpn_ssl_server

Provide a resource to create a VPN SSL Server.

## Example Usage

### Basic Configuration

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

### With Tags and DNS Configuration

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

### With SSO Authentication (Requires Whitelist)

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

### parameter only controls the feature switch. Detailed access policies must be configured through the TencentCloud console or other resources.

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

## Argument Reference

The following arguments are supported:

* `local_address` - (Required, List: [`String`]) List of local CIDR.
* `remote_address` - (Required, String) Remote CIDR for client.
* `ssl_vpn_server_name` - (Required, String) The name of ssl vpn server to be created.
* `vpn_gateway_id` - (Required, String, ForceNew) VPN gateway ID.
* `access_policy_enabled` - (Optional, Bool, ForceNew) Enable access policy control. Default: false.
* `compress` - (Optional, Bool) Need compressed. Currently is not supports compress. Default value: False.
* `dns_servers` - (Optional, List) DNS server configuration.
* `encrypt_algorithm` - (Optional, String) The encrypt algorithm. Valid values: AES-128-CBC, AES-192-CBC, AES-256-CBC.Default value: AES-128-CBC.
* `integrity_algorithm` - (Optional, String) The integrity algorithm. Valid values: SHA1. Default value: SHA1.
* `saml_data` - (Optional, String) SAML-DATA. Required when sso_enabled is true.
* `ssl_vpn_port` - (Optional, Int) The port of ssl vpn. Currently only supports UDP. Default value: 1194.
* `ssl_vpn_protocol` - (Optional, String) The protocol of ssl vpn. Default value: UDP.
* `sso_enabled` - (Optional, Bool) Enable SSO authentication. Default: false. This feature requires whitelist approval.
* `tags` - (Optional, Map) Tags for resource management.

The `dns_servers` object supports the following:

* `primary_dns` - (Optional, String) Primary DNS server address.
* `secondary_dns` - (Optional, String) Secondary DNS server address.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

VPN SSL Server can be imported, e.g.

```
$ terraform import tencentcloud_vpn_ssl_server.example vpns-cik6bjct
```

