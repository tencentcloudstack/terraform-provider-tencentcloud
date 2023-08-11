---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eni_sg_attachment"
sidebar_current: "docs-tencentcloud-resource-eni_sg_attachment"
description: |-
  Provides a resource to create a eni_sg_attachment
---

# tencentcloud_eni_sg_attachment

Provides a resource to create a eni_sg_attachment

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "vpc"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "example1" {
  name        = "tf-example-sg1"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

resource "tencentcloud_security_group" "example2" {
  name        = "tf-example-sg2"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

resource "tencentcloud_eni" "example" {
  name        = "tf-example-eni"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni desc."
  ipv4_count  = 1
}

resource "tencentcloud_eni_sg_attachment" "eni_sg_attachment" {
  network_interface_ids = [tencentcloud_eni.example.id]
  security_group_ids = [
    tencentcloud_security_group.example1.id,
    tencentcloud_security_group.example2.id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `network_interface_ids` - (Required, Set: [`String`], ForceNew) ENI instance ID. Such as:eni-pxir56ns. It Only support set one eni instance now.
* `security_group_ids` - (Required, Set: [`String`], ForceNew) Security group instance ID, for example:sg-33ocnj9n, can be obtained through DescribeSecurityGroups. There is a limit of 100 instances per request.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vpc eni_sg_attachment can be imported using the id, e.g.

```
terraform import tencentcloud_eni_sg_attachment.eni_sg_attachment eni_sg_attachment_id
```

