---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_classic_link_attachment"
sidebar_current: "docs-tencentcloud-resource-vpc_classic_link_attachment"
description: |-
  Provides a resource to create a vpc classic link attachment
---

# tencentcloud_vpc_classic_link_attachment

Provides a resource to create a vpc classic link attachment

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "cvm"
}

data "tencentcloud_images" "image" {
  image_type       = ["PUBLIC_IMAGE"]
  image_name_regex = "Final"
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
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_instance" "example" {
  instance_name            = "tf-example"
  availability_zone        = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  image_id                 = data.tencentcloud_images.image.images.0.image_id
  instance_type            = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type         = "CLOUD_PREMIUM"
  disable_security_service = true
  disable_monitor_service  = true
  vpc_id                   = tencentcloud_vpc.vpc.id
  subnet_id                = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_vpc_classic_link_attachment" "example" {
  vpc_id       = tencentcloud_vpc.vpc.id
  instance_ids = [tencentcloud_instance.example.id]
}
```

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Required, Set: [`String`], ForceNew) CVM instance ID. It only support set one instance now.
* `vpc_id` - (Required, String, ForceNew) VPC instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AttachClassicLinkVpc](https://www.tencentcloud.com/document/api/215/15779)
* Read: [DescribeClassicLinkInstances](https://www.tencentcloud.com/document/api/215/15776)
* Delete: [DetachClassicLinkVpc](https://www.tencentcloud.com/document/api/215/15777)

## Import

vpc classic link attachment can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_classic_link_attachment.example classic_link_attachment_id
```

