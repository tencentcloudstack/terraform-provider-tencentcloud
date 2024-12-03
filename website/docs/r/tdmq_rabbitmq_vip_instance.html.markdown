---
subcategory: "TDMQ for RabbitMQ(trabbit)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmq_rabbitmq_vip_instance"
sidebar_current: "docs-tencentcloud-resource-tdmq_rabbitmq_vip_instance"
description: |-
  Provides a resource to create a tdmq rabbitmq_vip_instance
---

# tencentcloud_tdmq_rabbitmq_vip_instance

Provides a resource to create a tdmq rabbitmq_vip_instance

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
```

## Argument Reference

The following arguments are supported:

* `cluster_name` - (Required, String) cluster name.
* `subnet_id` - (Required, String) Private network SubnetId.
* `vpc_id` - (Required, String) Private network VpcId.
* `zone_ids` - (Required, Set: [`Int`]) availability zone.
* `auto_renew_flag` - (Optional, Bool) Automatic renewal, the default is true.
* `enable_create_default_ha_mirror_queue` - (Optional, Bool) Mirrored queue, the default is false.
* `node_num` - (Optional, Int) The number of nodes, a minimum of 3 nodes for a multi-availability zone. If not passed, the default single availability zone is 1, and the multi-availability zone is 3.
* `node_spec` - (Optional, String) Node specifications. Valid values: rabbit-vip-basic-1 (for 4C8G), rabbit-vip-basic-2 (for 8C16G), rabbit-vip-basic-3 (for 12C24G), rabbit-vip-basic-4 (for 16C32G), rabbit-vip-basic-5 (for 2C4G), rabbit-vip-basic-6 (for 4C12G), rabbit-vip-basic-7 (for 8C24G). The default is rabbit-vip-basic-1. NOTE: The above specifications may be sold out or removed from the shelves.
* `storage_size` - (Optional, Int) Single node storage specification, the default is 200G.
* `time_span` - (Optional, Int) Purchase duration, the default is 1 (month).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

tdmq rabbitmq_vip_instance can be imported using the id, e.g.

```
terraform import tencentcloud_tdmq_rabbitmq_vip_instance.example amqp-mok52gmn
```

