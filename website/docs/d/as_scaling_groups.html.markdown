---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_scaling_groups"
sidebar_current: "docs-tencentcloud-datasource-as_scaling_groups"
description: |-
  Use this data source to query the detail information of an existing autoscaling group.
---

# tencentcloud_as_scaling_groups

Use this data source to query the detail information of an existing autoscaling group.

## Example Usage

```hcl
data "tencentcloud_as_scaling_groups" "as_scaling_groups" {
  scaling_group_name = "myasgroup"
  configuration_id   = "asc-oqio4yyj"
  result_output_file = "my_test_path"
}
```

## Argument Reference

The following arguments are supported:

* `configuration_id` - (Optional) Filter results by launch configuration ID.
* `result_output_file` - (Optional) Used to save results.
* `scaling_group_id` - (Optional) A specified scaling group ID used to query.
* `scaling_group_name` - (Optional) A scaling group name used to query.
* `tags` - (Optional) Tags used to query.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `scaling_group_list` - A list of scaling group. Each element contains the following attributes:
  * `configuration_id` - Launch configuration ID.
  * `create_time` - The time when the AS group was created.
  * `default_cooldown` - Default cooldown time of scaling group.
  * `desired_capacity` - The desired number of CVM instances.
  * `forward_balancer_ids` - A list of application clb ids.
    * `listener_id` - Listener ID for application load balancers.
    * `load_balancer_id` - ID of available load balancers.
    * `location_id` - ID of forwarding rules.
    * `target_attribute` - Attribute list of target rules.
      * `port` - Port number.
      * `weight` - Weight.
  * `instance_count` - Number of instance.
  * `load_balancer_ids` - A list of traditional clb ids which the CVM instances attached to.
  * `max_size` - The maximum number of CVM instances.
  * `min_size` - The minimum number of CVM instances.
  * `project_id` - ID of the project to which the scaling group belongs. Default value is 0.
  * `retry_policy` - A retry policy can be used when a creation fails.
  * `scaling_group_id` - Auto scaling group ID.
  * `scaling_group_name` - Auto scaling group name.
  * `status` - Current status of a scaling group.
  * `subnet_ids` - A list of subnet IDs.
  * `tags` - Tags of the scaling group.
  * `termination_policies` - A policy used to select a CVM instance to be terminated from the scaling group.
  * `vpc_id` - ID of the vpc with which the instance is associated.
  * `zones` - A list of available zones.


