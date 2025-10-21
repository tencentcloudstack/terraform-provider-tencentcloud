---
subcategory: "Gateway Load Balancer(GWLB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_gwlb_instance"
sidebar_current: "docs-tencentcloud-resource-gwlb_instance"
description: |-
  Provides a resource to create a gwlb gwlb_instance
---

# tencentcloud_gwlb_instance

Provides a resource to create a gwlb gwlb_instance

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
```

## Argument Reference

The following arguments are supported:

* `subnet_id` - (Required, String) Subnet ID of the VPC to which the backend target device of the GWLB belongs.
* `vpc_id` - (Required, String) ID of the VPC to which the backend target device of the GWLB belongs, such as vpc-12345678. It can be obtained through the DescribeVpcEx interface. If left blank, it defaults to DefaultVPC. This parameter is required when a private network CLB instance is created.
* `lb_charge_type` - (Optional, String) GWLB instance billing type, which currently supports POSTPAID_BY_HOUR only. The default is POSTPAID_BY_HOUR.
* `load_balancer_name` - (Optional, String) GWLB instance name. It supports input of 1 to 60 characters. If not filled in, it will be generated automatically by default.
* `tags` - (Optional, List) While the GWLB is purchased, it is tagged, with a maximum of 20 tag key-value pairs.

The `tags` object supports the following:

* `tag_key` - (Required, String) Tag key.
* `tag_value` - (Required, String) Tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time.
* `delete_protect` - Whether to turn on the deletion protection function.
* `isolated_time` - Time when the Gateway Load Balancer instance was isolated.
* `isolation` - 0: means not quarantined, 1: means quarantined.
* `operate_protect` - Whether to enable the configuration modification protection function.
* `status` - Gateway Load Balancer instance status. 0: Creating, 1: Running normally, 3: Removing.
* `target_group_id` - Unique ID of the associated target group.
* `vips` - Gateway Load Balancer provides virtual IP services.


## Import

gwlb gwlb_instance can be imported using the id, e.g.

```
terraform import tencentcloud_gwlb_instance.gwlb_instance gwlb_instance_id
```

