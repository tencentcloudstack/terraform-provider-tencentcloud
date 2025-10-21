---
subcategory: "TDMQ for RabbitMQ(trabbit)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_user"
sidebar_current: "docs-tencentcloud-resource-tdmq_rabbitmq_user"
description: |-
  Provides a resource to create a tdmq rabbitmq_user
---

# tencentcloud_tdmq_rabbitmq_user

Provides a resource to create a tdmq rabbitmq_user

## Example Usage

```hcl
data "tencentcloud_availability_zones" "zones" {
  name = "ap-guangzhou-6"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

# create vpc subnet
resource "tencentcloud_subnet" "subnet" {
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  availability_zone = "ap-guangzhou-6"
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

# create rabbitmq instance
resource "tencentcloud_tdmq_rabbitmq_vip_instance" "example" {
  zone_ids                              = [data.tencentcloud_availability_zones.zones.zones.0.id]
  vpc_id                                = tencentcloud_vpc.vpc.id
  subnet_id                             = tencentcloud_subnet.subnet.id
  cluster_name                          = "tf-example-rabbitmq-vip-instance"
  node_spec                             = "rabbit-vip-basic-1"
  node_num                              = 1
  storage_size                          = 200
  enable_create_default_ha_mirror_queue = false
  auto_renew_flag                       = true
  time_span                             = 1
}

# create rabbitmq user
resource "tencentcloud_tdmq_rabbitmq_user" "example" {
  instance_id     = tencentcloud_tdmq_rabbitmq_vip_instance.example.id
  user            = "tf-example-user"
  password        = "$Password"
  description     = "desc."
  tags            = ["management", "monitoring", "example"]
  max_connections = 3
  max_channels    = 3
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Cluster instance ID.
* `password` - (Required, String) Password, used when logging in.
* `user` - (Required, String) Username, used when logging in.
* `description` - (Optional, String) Describe.
* `max_channels` - (Optional, Int) The maximum number of channels for this user, if not filled in, there is no limit.
* `max_connections` - (Optional, Int) The maximum number of connections for this user, if not filled in, there is no limit.
* `tags` - (Optional, List: [`String`]) User tag, used to determine the permission range for changing user access to RabbitMQ Management. Management: regular console user, monitoring: management console user, other values: non console user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdmq rabbitmq_user can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rabbitmq_user.example amqp-8xzx822q#tf-example-user
```

