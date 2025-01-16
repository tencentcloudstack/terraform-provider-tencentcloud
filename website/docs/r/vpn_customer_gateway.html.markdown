---
subcategory: "VPN Connections(VPN)"
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
resource "tencentcloud_vpn_customer_gateway" "example" {
  name              = "tf-example"
  public_ip_address = "1.1.1.1"
  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the customer gateway. The length of character is limited to 1-60.
* `public_ip_address` - (Required, String, ForceNew) Public IP of the customer gateway.
* `bgp_asn` - (Optional, Int) BGP ASN. Value range: 1 - 4294967295. Using BGP requires configuring ASN. 139341, 45090, and 58835 are not available.
* `tags` - (Optional, Map) A list of tags used to associate different resources.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the customer gateway.


## Import

VPN customer gateway can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpn_customer_gateway.example cgw-xfqag
```

