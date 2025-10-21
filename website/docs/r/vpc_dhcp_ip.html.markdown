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
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_vpc_dhcp_ip" "example" {
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id
  dhcp_ip_name = "tf-example"
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

