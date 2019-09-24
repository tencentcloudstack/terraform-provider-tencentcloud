---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eip"
sidebar_current: "docs-tencentcloud-resource-eip"
description: |-
  Provides an EIP resource.
---

# tencentcloud_eip

Provides an EIP resource.

## Example Usage

```hcl
resource "tencentcloud_eip" "foo" {
  name = "awesome_gateway_ip"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of eip.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `public_ip` - The elastic ip address.
* `status` - The eip current status.


## Import

EIP can be imported using the id, e.g.

```
$ terraform import tencentcloud_eip.foo eip-nyvf60va
```

