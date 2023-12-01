---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_acl_attachment"
sidebar_current: "docs-tencentcloud-resource-vpc_acl_attachment"
description: |-
  Provide a resource to attach an existing subnet to Network ACL.
---

# tencentcloud_vpc_acl_attachment

Provide a resource to attach an existing subnet to Network ACL.

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
}

resource "tencentcloud_vpc_acl" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
  name   = "tf-example"
  ingress = [
    "ACCEPT#192.168.1.0/24#800#TCP",
    "ACCEPT#192.168.1.0/24#800-900#TCP",
  ]
  egress = [
    "ACCEPT#192.168.1.0/24#800#TCP",
    "ACCEPT#192.168.1.0/24#800-900#TCP",
  ]
}

resource "tencentcloud_vpc_acl_attachment" "attachment" {
  acl_id    = tencentcloud_vpc_acl.example.id
  subnet_id = tencentcloud_subnet.subnet.id
}
```

## Argument Reference

The following arguments are supported:

* `acl_id` - (Required, String, ForceNew) ID of the attached ACL.
* `subnet_id` - (Required, String, ForceNew) The Subnet instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Acl attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_vpc_acl_attachment.attachment acl-eotx5qsg#subnet-91x0geu6
```

