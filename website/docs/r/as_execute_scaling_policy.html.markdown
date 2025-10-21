---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_execute_scaling_policy"
sidebar_current: "docs-tencentcloud-resource-as_execute_scaling_policy"
description: |-
  Provides a resource to create a as execute_scaling_policy
---

# tencentcloud_as_execute_scaling_policy

Provides a resource to create a as execute_scaling_policy

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
  max_size           = 4
  min_size           = 1
  desired_capacity   = 2
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

resource "tencentcloud_as_execute_scaling_policy" "example" {
  auto_scaling_policy_id = tencentcloud_as_scaling_policy.example.id
  honor_cooldown         = false
  trigger_source         = "API"
}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_policy_id` - (Required, String, ForceNew) Auto-scaling policy ID. This parameter is not available to a target tracking policy.
* `honor_cooldown` - (Optional, Bool, ForceNew) Whether to check if the auto scaling group is in the cooldown period. Default value: false.
* `trigger_source` - (Optional, String, ForceNew) Source that triggers the scaling policy. Valid values: API and CLOUD_MONITOR. Default value: API. The value CLOUD_MONITOR is specific to the Cloud Monitor service.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

as execute_scaling_policy can be imported using the id, e.g.

```
terraform import tencentcloud_as_execute_scaling_policy.execute_scaling_policy execute_scaling_policy_id
```

