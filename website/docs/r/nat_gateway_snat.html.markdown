---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_nat_gateway_snat"
sidebar_current: "docs-tencentcloud-resource-nat_gateway_snat"
description: |-
  Provides a resource to create a NAT Gateway SNat rule.
---

# tencentcloud_nat_gateway_snat

Provides a resource to create a NAT Gateway SNat rule.

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "nat"
}

data "tencentcloud_images" "image" {
  os_name = "centos"
}

data "tencentcloud_instance_types" "instance_types" {
  filter {
    name   = "zone"
    values = [data.tencentcloud_availability_zones_by_product.zones.zones.0.name]
  }

  filter {
    name   = "instance-family"
    values = ["S5"]
  }

  cpu_core_count   = 2
  exclude_sold_out = true
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  route_table_id    = tencentcloud_route_table.route_table.id
}

resource "tencentcloud_eip" "eip_example1" {
  name = "eip_example1"
}

resource "tencentcloud_eip" "eip_example2" {
  name = "eip_example2"
}

# Create NAT Gateway
resource "tencentcloud_nat_gateway" "my_nat" {
  vpc_id         = tencentcloud_vpc.vpc.id
  name           = "tf_example_nat_gateway"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_example1.public_ip,
    tencentcloud_eip.eip_example2.public_ip,
  ]
}

# Create route_table and entry
resource "tencentcloud_route_table" "route_table" {
  vpc_id = tencentcloud_vpc.vpc.id
  name   = "tf_example"
}

resource "tencentcloud_route_table_entry" "route_entry" {
  route_table_id         = tencentcloud_route_table.route_table.id
  destination_cidr_block = "10.0.0.0/8"
  next_type              = "NAT"
  next_hub               = tencentcloud_nat_gateway.my_nat.id
}

# Subnet Nat gateway snat
resource "tencentcloud_nat_gateway_snat" "subnet_snat" {
  nat_gateway_id    = tencentcloud_nat_gateway.my_nat.id
  resource_type     = "SUBNET"
  subnet_id         = tencentcloud_subnet.subnet.id
  subnet_cidr_block = tencentcloud_subnet.subnet.cidr_block
  description       = "terraform test"
  public_ip_addr = [
    tencentcloud_eip.eip_example1.public_ip,
    tencentcloud_eip.eip_example2.public_ip,
  ]
}

# Create instance
resource "tencentcloud_instance" "example" {
  instance_name     = "tf_example"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  image_id          = data.tencentcloud_images.image.images.0.image_id
  instance_type     = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type  = "CLOUD_PREMIUM"
  system_disk_size  = 50
  hostname          = "user"
  project_id        = 0
  vpc_id            = tencentcloud_vpc.vpc.id
  subnet_id         = tencentcloud_subnet.subnet.id
}

# NetWorkInterface Nat gateway snat
resource "tencentcloud_nat_gateway_snat" "my_instance_snat" {
  nat_gateway_id           = tencentcloud_nat_gateway.my_nat.id
  resource_type            = "NETWORKINTERFACE"
  instance_id              = tencentcloud_instance.example.id
  instance_private_ip_addr = tencentcloud_instance.example.private_ip
  description              = "terraform test"
  public_ip_addr = [
    tencentcloud_eip.eip_example1.public_ip,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required, String) Description.
* `nat_gateway_id` - (Required, String, ForceNew) NAT gateway ID.
* `public_ip_addr` - (Required, List: [`String`]) Elastic IP address pool.
* `resource_type` - (Required, String, ForceNew) Resource type. Valid values: SUBNET, NETWORKINTERFACE.
* `instance_id` - (Optional, String, ForceNew) Instance ID, required when `resource_type` is NETWORKINTERFACE.
* `instance_private_ip_addr` - (Optional, String, ForceNew) Private IPs of the instance's primary ENI, required when `resource_type` is NETWORKINTERFACE.
* `subnet_cidr_block` - (Optional, String, ForceNew) The IPv4 CIDR of the subnet, required when `resource_type` is SUBNET.
* `subnet_id` - (Optional, String, ForceNew) Subnet instance ID, required when `resource_type` is SUBNET.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `snat_id` - SNAT rule ID.


## Import

VPN gateway route can be imported using the id, the id format must be '{nat_gateway_id}#{resource_id}', resource_id range `subnet_id`, `instance_id`, e.g.

SUBNET SNat
```
$ terraform import tencentcloud_nat_gateway_snat.my_snat nat-r4ip1cwt#subnet-2ap74y35
```

NETWORKINTERFACT SNat
```
$ terraform import tencentcloud_nat_gateway_snat.my_snat nat-r4ip1cwt#ins-da412f5a
```

