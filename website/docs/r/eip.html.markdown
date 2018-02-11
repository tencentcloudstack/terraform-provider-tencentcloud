---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip"
sidebar_current: "docs-tencentcloud-resource-eip"
description: |-
  Provides a TencentCloud EIP resource.
---

# tencentcloud_eip

Provides an EIP resource.

## Example Usage

Basic Usage

```hcl
resource "tencentcloud_eip" "foo" {
	name = "awesome_gateway_ip"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The eip's name. 


## Attributes Reference

The following attributes are exported:

* `id` - The EIP id, something like `eip-xxxxxxx`, use this for EIP assocication.
* `public_ip` - The elastic ip address.
* `status` - The EIP current status.