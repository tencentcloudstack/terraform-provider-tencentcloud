---
subcategory: "VPN Connections(VPN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_customer_gateway_configuration_download"
sidebar_current: "docs-tencentcloud-resource-vpn_customer_gateway_configuration_download"
description: |-
  Provides a resource to create a vpc vpn_customer_gateway_configuration_download
---

# tencentcloud_vpn_customer_gateway_configuration_download

Provides a resource to create a vpc vpn_customer_gateway_configuration_download

## Example Usage

```hcl
resource "tencentcloud_vpn_customer_gateway_configuration_download" "vpn_customer_gateway_configuration_download" {
  vpn_gateway_id    = "vpngw-gt8bianl"
  vpn_connection_id = "vpnx-kme2tx8m"
  customer_gateway_vendor {
    platform         = "comware"
    software_version = "V1.0"
    vendor_name      = "h3c"
  }
  interface_name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `customer_gateway_vendor` - (Required, List, ForceNew) Customer Gateway Vendor Info.
* `interface_name` - (Required, String, ForceNew) VPN connection access device physical interface name.
* `vpn_connection_id` - (Required, String, ForceNew) VPN Connection Instance id.
* `vpn_gateway_id` - (Required, String, ForceNew) VPN Gateway Instance ID.

The `customer_gateway_vendor` object supports the following:

* `platform` - (Required, String) Platform.
* `software_version` - (Required, String) SoftwareVersion.
* `vendor_name` - (Required, String) VendorName.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `customer_gateway_configuration` - xml configuration.


