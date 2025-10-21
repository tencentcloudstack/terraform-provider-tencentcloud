---
subcategory: "VPN Connections(VPN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_ssl_client"
sidebar_current: "docs-tencentcloud-resource-vpn_ssl_client"
description: |-
  Provide a resource to create a VPN SSL Client.
---

# tencentcloud_vpn_ssl_client

Provide a resource to create a VPN SSL Client.

## Example Usage

```hcl
resource "tencentcloud_vpn_ssl_client" "client" {
  ssl_vpn_server_id   = "vpns-aog5xcjj"
  ssl_vpn_client_name = "hello"
}
```

## Argument Reference

The following arguments are supported:

* `ssl_vpn_client_name` - (Required, String, ForceNew) The name of ssl vpn client to be created.
* `ssl_vpn_server_id` - (Required, String, ForceNew) VPN ssl server id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

VPN SSL Client can be imported, e.g.

```
$ terraform import tencentcloud_vpn_ssl_client.client vpn-client-id
```

