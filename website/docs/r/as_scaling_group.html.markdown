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

```hcl
resource "tencentcloud_as_scaling_group" "scaling_group" {
  scaling_group_name   = "tf-as-scaling-group"
  configuration_id     = "asc-oqio4yyj"
  max_size             = 1
  min_size             = 0
  vpc_id               = "vpc-3efmz0z"
  subnet_ids           = ["subnet-mc3egos"]
  project_id           = 0
  default_cooldown     = 400
  desired_capacity     = 1
  termination_policies = ["NEWEST_INSTANCE"]
  retry_policy         = "INCREMENTAL_INTERVALS"

  forward_balancer_ids {
    load_balancer_id = "lb-hk693b1l"
    listener_id      = "lbl-81wr497k"
    rule_id          = "loc-kiodx943"

    target_attribute {
      port   = 80
      weight = 90
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `configuration_id` - (Required) An available ID for a launch configuration.
* `max_size` - (Required) Maximum number of CVM instances. Valid value ranges: (0~2000).
* `min_size` - (Required) Minimum number of CVM instances. Valid value ranges: (0~2000).
* `scaling_group_name` - (Required) Name of a scaling group.
* `vpc_id` - (Required) ID of VPC network.
* `default_cooldown` - (Optional) Default cooldown time in second, and default value is `300`.
* `desired_capacity` - (Optional) Desired volume of CVM instances, which is between `max_size` and `min_size`.
* `forward_balancer_ids` - (Optional) List of application load balancers, which can't be specified with `load_balancer_ids` together.
* `load_balancer_ids` - (Optional) ID list of traditional load balancers.
* `multi_zone_subnet_policy` - (Optional) Multi zone or subnet strategy, Valid values: PRIORITY and EQUALITY.
* `project_id` - (Optional) Specifies to which project the scaling group belongs.
* `retry_policy` - (Optional) Available values for retry policies. Valid values: IMMEDIATE_RETRY and INCREMENTAL_INTERVALS.
* `subnet_ids` - (Optional) ID list of subnet, and for VPC it is required.
* `tags` - (Optional) Tags of a scaling group.
* `termination_policies` - (Optional) Available values for termination policies. Valid values: OLDEST_INSTANCE and NEWEST_INSTANCE.
* `zones` - (Optional) List of available zones, for Basic network it is required.

The `forward_balancer_ids` object supports the following:

* `listener_id` - (Required) Listener ID for application load balancers.
* `load_balancer_id` - (Required) ID of available load balancers.
* `target_attribute` - (Required) Attribute list of target rules.
* `rule_id` - (Optional) ID of forwarding rules.

The `target_attribute` object supports the following:

* `port` - (Required) Port number.
* `weight` - (Required) Weight.

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

