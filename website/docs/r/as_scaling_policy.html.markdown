---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_scaling_policy"
sidebar_current: "docs-tencentcloud-resource-as_scaling_policy"
description: |-
  Provides a resource for an AS (Auto scaling) policy.
---

# tencentcloud_as_scaling_policy

Provides a resource for an AS (Auto scaling) policy.

## Example Usage

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

resource "tencentcloud_as_scaling_policy" "example" {
  scaling_group_id    = tencentcloud_as_scaling_group.example.id
  policy_name         = "tf-as-scaling-policy"
  adjustment_type     = "EXACT_CAPACITY"
  adjustment_value    = 0
  comparison_operator = "GREATER_THAN"
  metric_name         = "CPU_UTILIZATION"
  threshold           = 80
  period              = 300
  continuous_time     = 10
  statistic           = "AVERAGE"
  cooldown            = 360
}
```

## Argument Reference

The following arguments are supported:

* `adjustment_type` - (Required, String) Specifies whether the adjustment is an absolute number or a percentage of the current capacity. Valid values: `CHANGE_IN_CAPACITY`, `EXACT_CAPACITY` and `PERCENT_CHANGE_IN_CAPACITY`.
* `adjustment_value` - (Required, Int) Define the number of instances by which to scale.For `CHANGE_IN_CAPACITY` type or PERCENT_CHANGE_IN_CAPACITY, a positive increment adds to the current capacity and a negative value removes from the current capacity. For `EXACT_CAPACITY` type, it defines an absolute number of the existing Auto Scaling group size.
* `comparison_operator` - (Required, String) Comparison operator. Valid values: `GREATER_THAN`, `GREATER_THAN_OR_EQUAL_TO`, `LESS_THAN`, `LESS_THAN_OR_EQUAL_TO`, `EQUAL_TO` and `NOT_EQUAL_TO`.
* `continuous_time` - (Required, Int) Retry times. Valid value ranges: (1~10).
* `metric_name` - (Required, String) Name of an indicator. Valid values: `CPU_UTILIZATION`, `MEM_UTILIZATION`, `LAN_TRAFFIC_OUT`, `LAN_TRAFFIC_IN`, `WAN_TRAFFIC_OUT` and `WAN_TRAFFIC_IN`.
* `period` - (Required, Int) Time period in second. Valid values: `60` and `300`.
* `policy_name` - (Required, String) Name of a policy used to define a reaction when an alarm is triggered.
* `scaling_group_id` - (Required, String, ForceNew) ID of a scaling group.
* `threshold` - (Required, Int) Alarm threshold.
* `cooldown` - (Optional, Int) Cooldwon time in second. Default is `30`0.
* `notification_user_group_ids` - (Optional, List: [`String`]) An ID group of users to be notified when an alarm is triggered.
* `statistic` - (Optional, String) Statistic types. Valid values: `AVERAGE`, `MAXIMUM` and `MINIMUM`. Default is `AVERAGE`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



