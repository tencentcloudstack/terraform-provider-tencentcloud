---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpn_customer_gateway"
sidebar_current: "docs-tencentcloud-resource-vpn_customer_gateway"
description: |-
  Provides a resource to create a VPN customer gateway.
---

# tencentcloud_vpn_customer_gateway

Provides a resource to create a VPN customer gateway.

## Example Usage

```hcl
resource "tencentcloud_vpn_customer_gateway" "foo" {
  name              = "test_vpn_customer_gateway"
  public_ip_address = "1.1.1.1"

  tags = {
    tag = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the customer gateway. The length of character is limited to 1-60.
* `public_ip_address` - (Required, ForceNew) Public ip of the customer gateway.
* `tags` - (Optional) A list of tags used to associate different resources.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `create_time` - Create time of the customer gateway.


## Import

VPN customer gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_customer_gateway.foo cgw-xfqag
```

