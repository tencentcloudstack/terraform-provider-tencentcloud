---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_subnet"
sidebar_current: "docs-tencentcloud-resource-subnet"
description: |-
  Provide a resource to create a VPC subnet.
---

# tencentcloud_subnet

Provide a resource to create a VPC subnet.

## Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  is_multicast      = false
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, String, ForceNew) The availability zone within which the subnet should be created.
* `cidr_block` - (Required, String, ForceNew) A network address block of the subnet.
* `name` - (Required, String) The name of subnet to be created.
* `vpc_id` - (Required, String, ForceNew) ID of the VPC to be associated.
* `is_multicast` - (Optional, Bool) Indicates whether multicast is enabled. The default value is 'true'.
* `route_table_id` - (Optional, String) ID of a routing table to which the subnet should be associated.
* `tags` - (Optional, Map) Tags of the subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `available_ip_count` - The number of available IPs.
* `create_time` - Creation time of subnet resource.
* `is_default` - Indicates whether it is the default VPC for this region.


## Import

Vpc subnet instance can be imported, e.g.

```
$ terraform import tencentcloud_subnet.test subnet_id
```

