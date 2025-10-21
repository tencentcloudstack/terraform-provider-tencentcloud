---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_schedule"
sidebar_current: "docs-tencentcloud-resource-as_schedule"
description: |-
  Provides a resource for an AS (Auto scaling) schedule.
---

# tencentcloud_as_schedule

Provides a resource for an AS (Auto scaling) schedule.

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

resource "tencentcloud_as_schedule" "example" {
  scaling_group_id     = tencentcloud_as_scaling_group.example.id
  schedule_action_name = "tf-as-schedule"
  max_size             = 10
  min_size             = 0
  desired_capacity     = 0
  start_time           = "2019-01-01T00:00:00+08:00"
  end_time             = "2019-12-01T00:00:00+08:00"
  recurrence           = "0 0 * * *"
}
```

## Argument Reference

The following arguments are supported:

* `desired_capacity` - (Required, Int) The desired number of CVM instances that should be running in the group.
* `max_size` - (Required, Int) The maximum size for the Auto Scaling group.
* `min_size` - (Required, Int) The minimum size for the Auto Scaling group.
* `scaling_group_id` - (Required, String, ForceNew) ID of a scaling group.
* `schedule_action_name` - (Required, String) The name of this scaling action.
* `start_time` - (Required, String) The time for this action to start, in "YYYY-MM-DDThh:mm:ss+08:00" format (UTC+8).
* `end_time` - (Optional, String) The time for this action to end, in "YYYY-MM-DDThh:mm:ss+08:00" format (UTC+8).
* `recurrence` - (Optional, String) The time when recurring future actions will start. Start time is specified by the user following the Unix cron syntax format. And this argument should be set with end_time together.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



