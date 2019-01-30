---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnat"
sidebar_current: "docs-tencentcloud-resource-vpc-dnat"
description: |-
  Provides a port mapping/forwarding of destination network address translation (DNAT/DNAPT) resource.
---

# tencentcloud_dnat

Provides a port mapping/forwarding of destination network address port translation (DNAT/DNAPT) resource.

## Example Usage

Basic usage:

```hcl
data "tencentcloud_availability_zones" "my_favorate_zones" {}

data "tencentcloud_image" "my_favorate_image" {
  filter {
    name   = "image-type"
    values = ["PUBLIC_IMAGE"]
  }
}

# Create VPC and Subnet
resource "tencentcloud_vpc" "main" {
  name       = "terraform test"
  cidr_block = "10.6.0.0/16"
}
resource "tencentcloud_subnet" "main_subnet" {
  vpc_id            = "${tencentcloud_vpc.main.id}"
  name              = "terraform test subnet"
  cidr_block        = "10.6.7.0/24"
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
}

# Create EIP
resource "tencentcloud_eip" "eip_dev_dnat" {
  name = "terraform_test"
}
resource "tencentcloud_eip" "eip_test_dnat" {
  name = "terraform_test"
}

# Create NAT Gateway
resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id           = "${tencentcloud_vpc.main.id}"
  name             = "terraform test"
  max_concurrent   = 3000000
  bandwidth        = 500
  assigned_eip_set = [
    "${tencentcloud_eip.eip_dev_dnat.public_ip}",
    "${tencentcloud_eip.eip_test_dnat.public_ip}",
  ]
}

# Create CVM
resource "tencentcloud_instance" "foo" {
  availability_zone = "${data.tencentcloud_availability_zones.my_favorate_zones.zones.0.name}"
  image_id          = "${data.tencentcloud_image.my_favorate_image.image_id}"
  vpc_id            = "${tencentcloud_vpc.main.id}"
  subnet_id         = "${tencentcloud_subnet.main_subnet.id}"
}

# Add DNAT Entry
resource "tencentcloud_dnat" "dev_dnat" {
  vpc_id       = "${tencentcloud_nat_gateway.my_nat.vpc_id}"
  nat_id       = "${tencentcloud_nat_gateway.my_nat.id}"
  protocol     = "tcp"
  elastic_ip   = "${tencentcloud_eip.eip_dev_dnat.public_ip}"
  elastic_port = "80"
  private_ip   = "${tencentcloud_instance.foo.private_ip}"
  private_port = "9001"
}
resource "tencentcloud_dnat" "test_dnat" {
  vpc_id       = "${tencentcloud_nat_gateway.my_nat.vpc_id}"
  nat_id       = "${tencentcloud_nat_gateway.my_nat.id}"
  protocol     = "udp"
  elastic_ip   = "${tencentcloud_eip.eip_test_dnat.public_ip}"
  elastic_port = "8080"
  private_ip   = "${tencentcloud_instance.foo.private_ip}"
  private_port = "9002"
}
```

## Argument Reference

The following arguments are supported:

* `nat_id` - (Required, Forces new resource) The ID for the NAT Gateway.
* `vpc_id` - (Required, Forces new resource) The VPC ID for the NAT Gateway.
* `protocol` - (Required, Forces new resource) The ip protocol, valid value is tcp|udp.
* `elastic_ip` - (Required, Forces new resource) The elastic IP of NAT gateway association, must a [Elastic IP](eip.html).
* `elastic_port` - (Required, Forces new resource) The external port, valid value is 1~65535.
* `private_ip` - (Required, Forces new resource) The internal ip, must a private ip (VPC IP).
* `private_port` (Required, Forces new resource) The internal port, valid value is 1~65535
