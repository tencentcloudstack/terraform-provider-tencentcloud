---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_reserve_ip_address"
sidebar_current: "docs-tencentcloud-resource-reserve_ip_address"
description: |-
  Provides a resource to create a vpc reserve ip addresses
---

# tencentcloud_reserve_ip_address

Provides a resource to create a vpc reserve ip addresses

## Example Usage

```hcl
resource "tencentcloud_reserve_ip_address" "reserve_ip" {
  vpc_id      = "xxxxxx"
  subnet_id   = "xxxxxx"
  ip_address  = "10.0.0.13"
  name        = "reserve-ip-tf"
  description = "description"
  tags = {
    "test1" = "test1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, String) VPC unique ID.
* `description` - (Optional, String) The IP description is retained on the intranet.
* `ip_address` - (Optional, String) Specify the reserved IP address of the intranet for which the IP application is requested.
* `name` - (Optional, String) The IP name is reserved for the intranet.
* `subnet_id` - (Optional, String) Subnet ID.
* `tags` - (Optional, Map) Tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `created_time` - Created time.
* `ip_type` - Ip type for product application.
* `reserve_ip_id` - Reserve ip ID.
* `resource_id` - The intranet retains the resource instance ID bound to the IPs.
* `state` - Binding status.


## Import

vpc reserve_ip_addresses can be imported using the id, e.g.

```
terraform import tencentcloud_reserve_ip_addresses.reserve_ip_addresses ${vpcId}#${reserveIpId}
```

