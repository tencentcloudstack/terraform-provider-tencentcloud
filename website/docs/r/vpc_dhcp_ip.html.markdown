---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_dhcp_ip"
sidebar_current: "docs-tencentcloud-resource-vpc_dhcp_ip"
description: |-
  Provides a resource to create a vpc dhcp_ip
---

# tencentcloud_vpc_dhcp_ip

Provides a resource to create a vpc dhcp_ip

## Example Usage

```hcl
resource "tencentcloud_vpc_dhcp_ip" "dhcp_ip" {
  vpc_id       = "vpc-1yg5ua6l"
  subnet_id    = "subnet-h7av55g8"
  dhcp_ip_name = "terraform-test"
}
```

## Argument Reference

The following arguments are supported:

* `dhcp_ip_name` - (Required, String) `DhcpIp` name.
* `subnet_id` - (Required, String) Subnet `ID`.
* `vpc_id` - (Required, String) The private network `ID`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc dhcp_ip can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_dhcp_ip.dhcp_ip dhcp_ip_id
```

