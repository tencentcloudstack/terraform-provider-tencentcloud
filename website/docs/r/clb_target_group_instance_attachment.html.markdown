---
subcategory: "Cloud Load Balancer(CLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clb_target_group_instance_attachment"
sidebar_current: "docs-tencentcloud-resource-clb_target_group_instance_attachment"
description: |-
  Provides a resource to create a CLB target group instance attachment.
---

# tencentcloud_clb_target_group_instance_attachment

Provides a resource to create a CLB target group instance attachment.

## Example Usage

```hcl
data "tencentcloud_availability_zones" "default" {}

data "tencentcloud_images" "images" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_instance_types" "instance_types" {
  cpu_core_count = 2
  memory_size    = 4
  filter {
    name   = "instance-family"
    values = ["S5"]
  }
}

resource "tencentcloud_vpc" "vpc" {
  cidr_block = "10.0.0.0/16"
  name       = "vpc"
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
  name              = "subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_instance" "example" {
  instance_name              = "tf-example"
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.images.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = tencentcloud_vpc.vpc.id
  subnet_id                  = tencentcloud_subnet.subnet.id
  internet_max_bandwidth_out = 100

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

data "tencentcloud_instances" "instances" {
  instance_id = tencentcloud_instance.example.id
}

resource "tencentcloud_clb_target_group" "example" {
  target_group_name = "tf-example"
  vpc_id            = tencentcloud_vpc.vpc.id
}

resource "tencentcloud_clb_target_group_instance_attachment" "example" {
  target_group_id = tencentcloud_clb_target_group.example.id
  bind_ip         = data.tencentcloud_instances.instances.instance_list[0].private_ip
  port            = 8080
  weight          = 10
}
```

## Argument Reference

The following arguments are supported:

* `bind_ip` - (Required, String, ForceNew) The Intranet IP of the target group instance.
* `port` - (Required, Int, ForceNew) Port of the target group instance.
* `target_group_id` - (Required, String, ForceNew) Target group ID.
* `weight` - (Required, Int) The weight of the target group instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CLB target group instance attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group_instance_attachment.example lbtg-3k3io0i0#172.16.48.18#8080
```

