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
data "tencentcloud_availability_zones" "my_zones" {}

data "tencentcloud_vpc" "my_vpc" {
  name = "Default-VPC"
}

data "tencentcloud_images" "my_image" {
  os_name = "centos"
}

data "tencentcloud_instance_types" "my_instance_types" {
  cpu_core_count = 1
  memory_size    = 1
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
  vpc_id         = data.tencentcloud_vpc.my_vpc.id
  name           = "terraform test"
  max_concurrent = 3000000
  bandwidth      = 500

  assigned_eip_set = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
}

# Create route_table and entry
resource "tencentcloud_route_table" "my_route_table" {
  vpc_id = data.tencentcloud_vpc.my_vpc.id
  name   = "terraform test"
}
resource "tencentcloud_route_table_entry" "my_route_entry" {
  route_table_id         = tencentcloud_route_table.my_route_table.id
  destination_cidr_block = "10.0.0.0/8"
  next_type              = "NAT"
  next_hub               = tencentcloud_nat_gateway.my_nat.id
}

# Create Subnet
resource "tencentcloud_subnet" "my_subnet" {
  vpc_id            = data.tencentcloud_vpc.my_vpc.id
  name              = "terraform test"
  cidr_block        = "172.29.23.0/24"
  availability_zone = data.tencentcloud_availability_zones.my_zones.zones.0.name
  route_table_id    = tencentcloud_route_table.my_route_table.id
}

# Subnet Nat gateway snat
resource "tencentcloud_nat_gateway_snat" "my_subnet_snat" {
  nat_gateway_id    = tencentcloud_nat_gateway.my_nat.id
  resource_type     = "SUBNET"
  subnet_id         = tencentcloud_subnet.my_subnet.id
  subnet_cidr_block = tencentcloud_subnet.my_subnet.cidr_block
  description       = "terraform test"
  public_ip_addr = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
    tencentcloud_eip.eip_test_dnat.public_ip,
  ]
}

# Create instance
resource "tencentcloud_instance" "my_instance" {
  instance_name              = "terraform test"
  availability_zone          = data.tencentcloud_availability_zones.my_zones.zones.0.name
  image_id                   = data.tencentcloud_images.my_image.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.my_instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = data.tencentcloud_vpc.my_vpc.id
  subnet_id                  = tencentcloud_subnet.my_subnet.id
  internet_max_bandwidth_out = 20
}

# NetWorkInterface Nat gateway snat
resource "tencentcloud_nat_gateway_snat" "my_instance_snat" {
  nat_gateway_id           = tencentcloud_nat_gateway.my_nat.id
  resource_type            = "NETWORKINTERFACE"
  instance_id              = tencentcloud_instance.my_instance.id
  instance_private_ip_addr = tencentcloud_instance.my_instance.private_ip
  description              = "terraform test"
  public_ip_addr = [
    tencentcloud_eip.eip_dev_dnat.public_ip,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Required) Description.
* `nat_gateway_id` - (Required, ForceNew) NAT gateway ID.
* `public_ip_addr` - (Required) Elastic IP address pool.
* `resource_type` - (Required, ForceNew) Resource type. Valid values: SUBNET, NETWORKINTERFACE.
* `instance_id` - (Optional, ForceNew) Instance ID, required when `resource_type` is NETWORKINTERFACE.
* `instance_private_ip_addr` - (Optional, ForceNew) Private IPs of the instance's primary ENI, required when `resource_type` is NETWORKINTERFACE.
* `subnet_cidr_block` - (Optional, ForceNew) The IPv4 CIDR of the subnet, required when `resource_type` is SUBNET.
* `subnet_id` - (Optional, ForceNew) Subnet instance ID, required when `resource_type` is SUBNET.

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

