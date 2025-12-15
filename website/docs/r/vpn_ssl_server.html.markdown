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

## Argument Reference

The following arguments are supported:

* `local_address` - (Required, List: [`String`]) List of local CIDR.
* `remote_address` - (Required, String) Remote CIDR for client.
* `ssl_vpn_server_name` - (Required, String) The name of ssl vpn server to be created.
* `vpn_gateway_id` - (Required, String, ForceNew) VPN gateway ID.
* `compress` - (Optional, Bool) Need compressed. Currently is not supports compress. Default value: False.
* `encrypt_algorithm` - (Optional, String) The encrypt algorithm. Valid values: AES-128-CBC, AES-192-CBC, AES-256-CBC.Default value: AES-128-CBC.
* `integrity_algorithm` - (Optional, String) The integrity algorithm. Valid values: SHA1. Default value: SHA1.
* `ssl_vpn_port` - (Optional, Int) The port of ssl vpn. Currently only supports UDP. Default value: 1194.
* `ssl_vpn_protocol` - (Optional, String) The protocol of ssl vpn. Default value: UDP.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

VPN SSL Server can be imported, e.g.

```
$ terraform import tencentcloud_vpn_ssl_server.example vpns-cik6bjct
```

