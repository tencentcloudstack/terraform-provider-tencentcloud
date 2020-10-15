---
subcategory: "CLB"
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
data "tencentcloud_images" "my_favorite_image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "centos"
}

data "tencentcloud_instance_types" "my_favorite_instance_types" {
  filter {
    name   = "instance-family"
    values = ["S3"]
  }

  cpu_core_count = 1
  memory_size    = 1
}

data "tencentcloud_availability_zones" "default" {
}

resource "tencentcloud_vpc" "app" {
  cidr_block = "10.0.0.0/16"
  name       = "awesome_app_vpc"
}

resource "tencentcloud_subnet" "app" {
  vpc_id            = tencentcloud_vpc.app.id
  availability_zone = data.tencentcloud_availability_zones.default.zones.0.name
  name              = "awesome_app_subnet"
  cidr_block        = "10.0.1.0/24"
}

resource "tencentcloud_instance" "my_awesome_app" {
  instance_name              = "awesome_app"
  availability_zone          = data.tencentcloud_availability_zones.default.zones.0.name
  image_id                   = data.tencentcloud_images.my_favorite_image.images.0.image_id
  instance_type              = data.tencentcloud_instance_types.my_favorite_instance_types.instance_types.0.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  hostname                   = "user"
  project_id                 = 0
  vpc_id                     = tencentcloud_vpc.app.id
  subnet_id                  = tencentcloud_subnet.app.id
  internet_max_bandwidth_out = 20

  data_disks {
    data_disk_type = "CLOUD_PREMIUM"
    data_disk_size = 50
    encrypt        = false
  }

  tags = {
    tagKey = "tagValue"
  }
}

data "tencentcloud_instances" "foo" {
  instance_id = tencentcloud_instance.my_awesome_app.id
}

resource "tencentcloud_clb_target_group" "test" {
  target_group_name = "test"
  vpc_id            = tencentcloud_vpc.app.id
  port              = 33
}

resource "tencentcloud_clb_target_group_instance_attachment" "test" {
  target_group_id = tencentcloud_clb_targetgroup.test.id
  bind_ip         = data.tencentcloud_instances.foo.instance_list[0].private_ip
  port            = 222
  weight          = 3
}
```

## Argument Reference

The following arguments are supported:

* `bind_ip` - (Required, ForceNew) The Intranet IP of the target group instance.
* `port` - (Required, ForceNew) Port of the target group instance.
* `target_group_id` - (Required, ForceNew) Target group ID.
* `weight` - (Required) The weight of the target group instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

CLB target group instance attachment can be imported using the id, e.g.

```
$ terraform import tencentcloud_clb_target_group_instance_attachment.test lbtg-3k3io0i0#172.16.48.18#222
```

