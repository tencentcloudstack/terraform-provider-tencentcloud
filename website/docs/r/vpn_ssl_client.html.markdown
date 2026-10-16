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

### Basic Configuration

```hcl
resource "tencentcloud_vpn_ssl_client" "client" {
  ssl_vpn_server_id   = "vpns-aog5xcjj"
  ssl_vpn_client_name = "hello"
}
```

### With Tags

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

## Argument Reference

The following arguments are supported:

* `ssl_vpn_client_name` - (Required, String, ForceNew) The name of ssl vpn client to be created.
* `ssl_vpn_server_id` - (Required, String, ForceNew) VPN ssl server id.
* `tags` - (Optional, Map) Tags of the VPN SSL client.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateVpnGatewaySslClient](https://cloud.tencent.com/document/api/215/70290)
* Create: [DescribeVpcTaskResult](https://cloud.tencent.com/document/api/215/59037)
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319)
* Read: [DescribeVpnGatewaySslClients](https://cloud.tencent.com/document/api/215/70286)
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322)
* Delete: [DeleteVpnGatewaySslClient](https://cloud.tencent.com/document/api/215/70288)
* Delete: [DescribeVpcTaskResult](https://cloud.tencent.com/document/api/215/59037)

## Import

VPN SSL Client can be imported, e.g.

```
$ terraform import tencentcloud_vpn_ssl_client.client vpn-client-id
```

