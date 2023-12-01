---
subcategory: "VPN Connections(VPN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_gateway_ssl_client_cert"
sidebar_current: "docs-tencentcloud-resource-vpn_gateway_ssl_client_cert"
description: |-
  Provides a resource to create a vpc vpn_gateway_ssl_client_cert
---

# tencentcloud_vpn_gateway_ssl_client_cert

Provides a resource to create a vpc vpn_gateway_ssl_client_cert

## Example Usage

```hcl
resource "tencentcloud_vpn_gateway_ssl_client_cert" "vpn_gateway_ssl_client_cert" {
  ssl_vpn_client_id = "vpnc-123456"
  switch            = "off"
}
```

## Argument Reference

The following arguments are supported:

* `ssl_vpn_client_id` - (Required, String) SSL-VPN-CLIENT Instance ID.
* `switch` - (Optional, String) `on`: Enable, `off`: Disable.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc vpn_gateway_ssl_client_cert can be imported using the id, e.g.

```
terraform import tencentcloud_vpn_gateway_ssl_client_cert.vpn_gateway_ssl_client_cert ssl_client_id
```

