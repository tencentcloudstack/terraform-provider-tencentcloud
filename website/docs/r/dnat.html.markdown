---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dnat"
sidebar_current: "docs-tencentcloud-resource-dnat"
description: |-
  Provides a resource to create a NAT forwarding.
---

# tencentcloud_dnat

Provides a resource to create a NAT forwarding.

## Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {}

data "tencentcloud_images" "example" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
}

data "tencentcloud_instance_types" "instance_types" {
  filter {
    name   = "zone"
    values = [data.tencentcloud_availability_zones.zones.zones.0.name]
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
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  name              = "example-vpc"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_eip" "eip_example" {
  name = "tf_nat_gateway_eip"
}

resource "tencentcloud_nat_gateway" "example" {
  name           = "tf_example_nat_gateway"
  vpc_id         = tencentcloud_vpc.vpc.id
  bandwidth      = 100
  max_concurrent = 1000000
  assigned_eip_set = [
    tencentcloud_eip.eip_example.public_ip,
  ]
  tags = {
    tf_tag_key = "tf_tag_value"
  }
}

resource "tencentcloud_instance" "example" {
  instance_name              = "tf_example_instance"
  availability_zone          = data.tencentcloud_availability_zones.zones.zones.0.name
  image_id                   = data.tencentcloud_images.example.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  allocate_public_ip         = true
  internet_max_bandwidth_out = 10
  vpc_id                     = tencentcloud_vpc.vpc.id
  subnet_id                  = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_dnat" "example" {
  vpc_id       = tencentcloud_vpc.vpc.id
  nat_id       = tencentcloud_nat_gateway.example.id
  protocol     = "TCP"
  elastic_ip   = tencentcloud_eip.eip_example.public_ip
  elastic_port = 80
  private_ip   = tencentcloud_instance.example.private_ip
  private_port = 9090
  description  = "desc."
}
```

## Argument Reference

The following arguments are supported:

* `elastic_ip` - (Required, String, ForceNew) Network address of the EIP.
* `elastic_port` - (Required, String, ForceNew) Port of the EIP.
* `nat_id` - (Required, String, ForceNew) ID of the NAT gateway.
* `private_ip` - (Required, String, ForceNew) Network address of the backend service.
* `private_port` - (Required, String, ForceNew) Port of intranet.
* `protocol` - (Required, String, ForceNew) Type of the network protocol. Valid value: `TCP` and `UDP`.
* `vpc_id` - (Required, String, ForceNew) ID of the VPC.
* `description` - (Optional, String) Description of the NAT forward.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

NAT forwarding can be imported using the id, e.g.

```
$ terraform import tencentcloud_dnat.foo tcp://vpc-asg3sfa3:nat-1asg3t63@127.15.2.3:8080
```

