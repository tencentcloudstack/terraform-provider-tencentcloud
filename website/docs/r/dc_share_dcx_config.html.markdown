---
subcategory: "Direct Connect(DC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_share_dcx_config"
sidebar_current: "docs-tencentcloud-resource-dc_share_dcx_config"
description: |-
  Provides a resource to create a dc share dcx config
---

# tencentcloud_dc_share_dcx_config

Provides a resource to create a dc share dcx config

## Example Usage

```hcl
resource "tencentcloud_dc_share_dcx_config" "example" {
  direct_connect_tunnel_id = "dcx-4z49tnws"
  enable                   = true
}
```

## Argument Reference

The following arguments are supported:

* `direct_connect_tunnel_id` - (Required, String) the direct connect owner accept or reject the apply of direct connect tunnel.
* `enable` - (Required, Bool) if accept or reject direct connect tunnel.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AcceptDirectConnectTunnel](https://cloud.tencent.com/document/api/216/19822) - This API is used to accept an application for a dedicated tunnel.
* Create: [RejectDirectConnectTunnel](https://cloud.tencent.com/document/api/216/19817) - This API is used to reject an application for a dedicated tunnel.
* Read: [DescribeDirectConnectTunnels](https://cloud.tencent.com/document/api/216/19819) - This API is used to query the list of dedicated tunnels.
* Update: [AcceptDirectConnectTunnel](https://cloud.tencent.com/document/api/216/19822) - This API is used to accept an application for a dedicated tunnel.
* Update: [RejectDirectConnectTunnel](https://cloud.tencent.com/document/api/216/19817) - This API is used to reject an application for a dedicated tunnel.

## Import

dc share dcx config can be imported using the id, e.g.

```
terraform import tencentcloud_dc_share_dcx_config.example dcx-4z49tnws
```

