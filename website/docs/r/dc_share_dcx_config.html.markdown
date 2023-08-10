---
subcategory: "Direct Connect(DC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dc_share_dcx_config"
sidebar_current: "docs-tencentcloud-resource-dc_share_dcx_config"
description: |-
  Provides a resource to create a dc share_dcx_config
---

# tencentcloud_dc_share_dcx_config

Provides a resource to create a dc share_dcx_config

## Example Usage

```hcl
resource "tencentcloud_dc_share_dcx_config" "share_dcx_config" {
  direct_connect_tunnel_id = "dcx-4z49tnws"
  enable                   = false
}
```

## Argument Reference

The following arguments are supported:

* `direct_connect_tunnel_id` - (Required, String) the direct connect owner accept or reject the apply of direct connect tunnel.
* `enable` - (Required, Bool) if accept or reject direct connect tunnel.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dc share_dcx_config can be imported using the id, e.g.

```
terraform import tencentcloud_dc_share_dcx_config.share_dcx_config dcx_id
```

