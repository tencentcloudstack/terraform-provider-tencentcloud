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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeVpcTaskResult](https://www.tencentcloud.com/document/product/215/41567) - This API is used to query the execution result of a VPC task.
* Create: [DisableVpnGatewaySslClientCert](https://cloud.tencent.com/document/api/215/70284)
* Create: [EnableVpnGatewaySslClientCert](https://cloud.tencent.com/document/api/215/70282)
* Read: [DescribeVpnGatewaySslClients](https://cloud.tencent.com/document/api/215/70286)
* Update: [DescribeVpcTaskResult](https://www.tencentcloud.com/document/product/215/41567) - This API is used to query the execution result of a VPC task.
* Update: [DisableVpnGatewaySslClientCert](https://cloud.tencent.com/document/api/215/70284)
* Update: [EnableVpnGatewaySslClientCert](https://cloud.tencent.com/document/api/215/70282)

## Import

vpc vpn_gateway_ssl_client_cert can be imported using the id, e.g.

```
terraform import tencentcloud_vpn_gateway_ssl_client_cert.vpn_gateway_ssl_client_cert ssl_client_id
```

