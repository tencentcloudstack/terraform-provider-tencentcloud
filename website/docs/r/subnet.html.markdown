---
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
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "foo" {
  name       = "guagua-ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = "${var.availability_zone}"
  name              = "guagua-ci-temp-test"
  vpc_id            = "${tencentcloud_vpc.foo.id}"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, ForceNew) The availability zone within which the subnet should be created.
* `cidr_block` - (Required, ForceNew) A network address block of the subnet.
* `name` - (Required) The name of subnet to be created.
* `vpc_id` - (Required, ForceNew) ID of the VPC to be associated.
* `is_multicast` - (Optional) Indicates whether multicast is enabled. The default value is 'true'.
* `route_table_id` - (Optional) ID of a routing table to which the subnet should be associated.
* `tags` - (Optional) Tags of the subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `available_ip_count` - The number of available IPs.
* `create_time` - Creation time of subnet resource.
* `is_default` - Indicates whether it is the default VPC for this region.


## Import

Vpc subnet instance can be imported, e.g.

```
$ terraform import tencentcloud_subnet.test subnet_id
```

