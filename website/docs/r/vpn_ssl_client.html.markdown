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
* Create: [DescribeVpcTaskResult](https://www.tencentcloud.com/document/product/215/41567) - This API is used to query the execution result of a VPC task.
* Read: [DescribeResourceTagsByResourceIds](https://www.tencentcloud.com/document/product/651/33022) - This API is used to query the tag key-value pairs associated with an existing resource.
* Read: [DescribeVpnGatewaySslClients](https://cloud.tencent.com/document/api/215/70286)
* Update: [ModifyResourceTags](https://www.tencentcloud.com/document/product/651/33012) - This API is used to modify all tags associated with a resource.
* Delete: [DeleteVpnGatewaySslClient](https://cloud.tencent.com/document/api/215/70288)
* Delete: [DescribeVpcTaskResult](https://www.tencentcloud.com/document/product/215/41567) - This API is used to query the execution result of a VPC task.

## Import

VPN SSL Client can be imported, e.g.

```
$ terraform import tencentcloud_vpn_ssl_client.client vpn-client-id
```

