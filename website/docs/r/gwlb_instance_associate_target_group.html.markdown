---
subcategory: "Gateway Load Balancer(GWLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gwlb_instance_associate_target_group"
sidebar_current: "docs-tencentcloud-resource-gwlb_instance_associate_target_group"
description: |-
  Provides a resource to create a gwlb gwlb_instance_associate_target_groups
---

# tencentcloud_gwlb_instance_associate_target_group

Provides a resource to create a gwlb gwlb_instance_associate_target_groups

## Example Usage

```hcl
variable "availability_zone" {
  default = "ap-guangzhou-3"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.20.0/28"
  is_multicast      = false
}

resource "tencentcloud_gwlb_instance" "gwlb_instance" {
  vpc_id             = tencentcloud_vpc.vpc.id
  subnet_id          = tencentcloud_subnet.subnet.id
  load_balancer_name = "tf-test"
  lb_charge_type     = "POSTPAID_BY_HOUR"
  tags {
    tag_key   = "test_key"
    tag_value = "tag_value"
  }
}

resource "tencentcloud_gwlb_target_group" "gwlb_target_group" {
  target_group_name = "tf-test"
  vpc_id            = tencentcloud_vpc.vpc.id
  port              = 6081
  health_check {
    health_switch = true
    protocol      = "tcp"
    port          = 6081
    timeout       = 2
    interval_time = 5
    health_num    = 3
    un_health_num = 3
  }
}

resource "tencentcloud_gwlb_instance_associate_target_group" "gwlb_instance_associate_target_group" {
  load_balancer_id = tencentcloud_gwlb_instance.gwlb_instance.id
  target_group_id  = tencentcloud_gwlb_target_group.gwlb_target_group.id
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required, String, ForceNew) GWLB instance ID.
* `target_group_id` - (Required, String, ForceNew) Target group ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



