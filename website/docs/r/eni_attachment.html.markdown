---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eni_attachment"
sidebar_current: "docs-tencentcloud-resource-eni_attachment"
description: |-
  Provides a resource to detailed information of attached backend server to an ENI.
---

# tencentcloud_eni_attachment

Provides a resource to detailed information of attached backend server to an ENI.

## Example Usage

```hcl
data "tencentcloud_images" "images" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_instance_types" "instance_types" {
  filter {
    name   = "instance-family"
    values = ["S3"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name       = "ci-test-eni-vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = "ap-guangzhou-3"
  name              = "ci-test-eni-subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_eni" "eni" {
  name        = "ci-test-eni"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni desc"
  ipv4_count  = 1
}

resource "tencentcloud_instance" "example" {
  instance_name            = "ci-test-eni-attach"
  availability_zone        = data.tencentcloud_availability_zones.zones.zones.0.name
  image_id                 = data.tencentcloud_images.images.images.0.image_id
  instance_type            = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type         = "CLOUD_PREMIUM"
  disable_security_service = true
  disable_monitor_service  = true
  vpc_id                   = tencentcloud_vpc.vpc.id
  subnet_id                = tencentcloud_subnet.subnet.id
}

resource "tencentcloud_eni_attachment" "example" {
  eni_id      = tencentcloud_eni.eni.id
  instance_id = tencentcloud_instance.example.id
}
```

## Argument Reference

The following arguments are supported:

* `eni_id` - (Required, String, ForceNew) ID of the ENI.
* `instance_id` - (Required, String, ForceNew) ID of the instance which bind the ENI.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ENI attachment can be imported using the id, e.g.

```
terraform import tencentcloud_eni_attachment.example eni-gtlvkjvz+ins-0h3a5new
```

