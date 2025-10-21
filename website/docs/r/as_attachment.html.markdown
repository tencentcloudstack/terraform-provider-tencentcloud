---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_attachment"
sidebar_current: "docs-tencentcloud-resource-as_attachment"
description: |-
  Provides a resource to attach or detach CVM instances to a specified scaling group.
---

# tencentcloud_as_attachment

Provides a resource to attach or detach CVM instances to a specified scaling group.

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "as"
}

data "tencentcloud_images" "image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
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
}

resource "tencentcloud_as_scaling_config" "example" {
  configuration_name = "tf-example"
  image_id           = data.tencentcloud_images.image.images.0.image_id
  instance_types     = ["SA1.SMALL1", "SA2.SMALL1", "SA2.SMALL2", "SA2.SMALL4"]
  instance_name_settings {
    instance_name = "test-ins-name"
  }
}

resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name = "tf-example"
  configuration_id   = tencentcloud_as_scaling_config.example.id
  max_size           = 1
  min_size           = 0
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_ids         = [tencentcloud_subnet.subnet.id]
}

resource "tencentcloud_instance" "example" {
  instance_name              = "tf_example_instance"
  availability_zone          = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  image_id                   = data.tencentcloud_images.image.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  allocate_public_ip         = true
  internet_max_bandwidth_out = 10
  vpc_id                     = tencentcloud_vpc.vpc.id
  subnet_id                  = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_as_attachment" "attachment" {
  scaling_group_id = tencentcloud_as_scaling_group.example.id
  instance_ids     = [tencentcloud_instance.example.id]
}
```

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Required, Set: [`String`]) ID list of CVM instances to be attached to the scaling group.
* `scaling_group_id` - (Required, String, ForceNew) ID of a scaling group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



