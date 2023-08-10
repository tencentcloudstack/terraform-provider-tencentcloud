---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_scaling_group"
sidebar_current: "docs-tencentcloud-resource-as_scaling_group"
description: |-
  Provides a resource to create a group of AS (Auto scaling) instances.
---

# tencentcloud_as_scaling_group

Provides a resource to create a group of AS (Auto scaling) instances.

## Example Usage

### Create a basic Scaling Group

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "as"
}

data "tencentcloud_images" "image" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "TencentOS Server 3.2 (Final)"
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
```

### Create a complete Scaling Group

```hcl
resource "tencentcloud_clb_instance" "example" {
  network_type = "INTERNAL"
  clb_name     = "clb-example"
  project_id   = 0
  vpc_id       = tencentcloud_vpc.vpc.id
  subnet_id    = tencentcloud_subnet.subnet.id

  tags = {
    test = "tf"
  }
}

resource "tencentcloud_clb_listener" "example" {
  clb_id        = tencentcloud_clb_instance.example.id
  listener_name = "listener-example"
  port          = 80
  protocol      = "HTTP"
}

resource "tencentcloud_clb_listener_rule" "example" {
  listener_id = tencentcloud_clb_listener.example.listener_id
  clb_id      = tencentcloud_clb_instance.example.id
  domain      = "foo.net"
  url         = "/bar"
}

resource "tencentcloud_as_scaling_group" "example" {
  scaling_group_name   = "tf-example"
  configuration_id     = tencentcloud_as_scaling_config.example.id
  max_size             = 1
  min_size             = 0
  vpc_id               = tencentcloud_vpc.vpc.id
  subnet_ids           = [tencentcloud_subnet.subnet.id]
  project_id           = 0
  default_cooldown     = 400
  desired_capacity     = 1
  termination_policies = ["NEWEST_INSTANCE"]
  retry_policy         = "INCREMENTAL_INTERVALS"

  forward_balancer_ids {
    load_balancer_id = tencentcloud_clb_instance.example.id
    listener_id      = tencentcloud_clb_listener.example.listener_id
    rule_id          = tencentcloud_clb_listener_rule.example.rule_id

    target_attribute {
      port   = 80
      weight = 90
    }
  }

  tags = {
    "createBy" = "tfExample"
  }
}
```

## Argument Reference

The following arguments are supported:

* `configuration_id` - (Required, String) An available ID for a launch configuration.
* `max_size` - (Required, Int) Maximum number of CVM instances. Valid value ranges: (0~2000).
* `min_size` - (Required, Int) Minimum number of CVM instances. Valid value ranges: (0~2000).
* `scaling_group_name` - (Required, String) Name of a scaling group.
* `vpc_id` - (Required, String) ID of VPC network.
* `default_cooldown` - (Optional, Int) Default cooldown time in second, and default value is `300`.
* `desired_capacity` - (Optional, Int) Desired volume of CVM instances, which is between `max_size` and `min_size`.
* `forward_balancer_ids` - (Optional, List) List of application load balancers, which can't be specified with `load_balancer_ids` together.
* `load_balancer_ids` - (Optional, List: [`String`]) ID list of traditional load balancers.
* `multi_zone_subnet_policy` - (Optional, String) Multi zone or subnet strategy, Valid values: PRIORITY and EQUALITY.
* `project_id` - (Optional, Int) Specifies to which project the scaling group belongs.
* `replace_load_balancer_unhealthy` - (Optional, Bool) Enable unhealthy instance replacement. If set to `true`, AS will replace instances that are found unhealthy in the CLB health check.
* `replace_monitor_unhealthy` - (Optional, Bool) Enables unhealthy instance replacement. If set to `true`, AS will replace instances that are flagged as unhealthy by Cloud Monitor.
* `retry_policy` - (Optional, String) Available values for retry policies. Valid values: IMMEDIATE_RETRY and INCREMENTAL_INTERVALS.
* `scaling_mode` - (Optional, String) Indicates scaling mode which creates and terminates instances (classic method), or method first tries to start stopped instances (wake up stopped) to perform scaling operations. Available values: `CLASSIC_SCALING`, `WAKE_UP_STOPPED_SCALING`. Default: `CLASSIC_SCALING`.
* `subnet_ids` - (Optional, List: [`String`]) ID list of subnet, and for VPC it is required.
* `tags` - (Optional, Map) Tags of a scaling group.
* `termination_policies` - (Optional, List: [`String`]) Available values for termination policies. Valid values: OLDEST_INSTANCE and NEWEST_INSTANCE.
* `zones` - (Optional, List: [`String`]) List of available zones, for Basic network it is required.

The `forward_balancer_ids` object supports the following:

* `listener_id` - (Required, String) Listener ID for application load balancers.
* `load_balancer_id` - (Required, String) ID of available load balancers.
* `target_attribute` - (Required, List) Attribute list of target rules.
* `rule_id` - (Optional, String) ID of forwarding rules.

The `target_attribute` object supports the following:

* `port` - (Required, Int) Port number.
* `weight` - (Required, Int) Weight.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time when the AS group was created.
* `instance_count` - Instance number of a scaling group.
* `status` - Current status of a scaling group.


## Import

AutoScaling Groups can be imported using the id, e.g.

```
$ terraform import tencentcloud_as_scaling_group.scaling_group asg-n32ymck2
```

