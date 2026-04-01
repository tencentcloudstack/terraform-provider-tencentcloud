---
subcategory: "Private Link(PLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_end_point"
sidebar_current: "docs-tencentcloud-resource-vpc_end_point"
description: |-
  Provides a resource to create a VPC end point
---

# tencentcloud_vpc_end_point

Provides a resource to create a VPC end point

## Example Usage

### Basic

```hcl
resource "tencentcloud_vpc_end_point" "example" {
  vpc_id               = "vpc-391sv4w3"
  subnet_id            = "subnet-ljyn7h30"
  end_point_name       = "tf-example"
  end_point_service_id = "vpcsvc-69y13tdb"
  end_point_vip        = "10.0.2.1"

  security_groups_ids = [
    "sg-ghvp9djf",
    "sg-if748odn",
    "sg-3k7vtgf7",
  ]
}
```

### With SecurityGroupId

```hcl
resource "tencentcloud_vpc_end_point" "example" {
  vpc_id               = "vpc-391sv4w3"
  subnet_id            = "subnet-ljyn7h30"
  end_point_name       = "tf-example"
  end_point_service_id = "vpcsvc-69y13tdb"
  end_point_vip        = "10.0.2.1"
  security_group_id    = "sg-ghvp9djf"
}
```

### With Tags

```hcl
resource "tencentcloud_vpc_end_point" "example" {
  vpc_id               = "vpc-391sv4w3"
  subnet_id            = "subnet-ljyn7h30"
  end_point_name       = "tf-example"
  end_point_service_id = "vpcsvc-69y13tdb"
  end_point_vip        = "10.0.2.1"

  tags {
    key   = "env"
    value = "prod"
  }
  tags {
    key   = "owner"
    value = "terraform"
  }
}
```

### With IpAddressType

```hcl
resource "tencentcloud_vpc_end_point" "example" {
  vpc_id               = "vpc-391sv4w3"
  subnet_id            = "subnet-ljyn7h30"
  end_point_name       = "tf-example"
  end_point_service_id = "vpcsvc-69y13tdb"
  end_point_vip        = "10.0.2.1"
  ip_address_type      = "Ipv4"
}
```

## Argument Reference

The following arguments are supported:

* `end_point_name` - (Required, String) Name of endpoint.
* `end_point_service_id` - (Required, String) ID of endpoint service.
* `subnet_id` - (Required, String) ID of subnet instance.
* `vpc_id` - (Required, String) ID of vpc instance.
* `end_point_vip` - (Optional, String) VIP of endpoint ip.
* `security_groups_ids` - (Optional, List: [`String`]) Ordered security groups associated with the endpoint.
* `security_group_id` - (Optional, String) ID of security group.
* `tags` - (Optional, List) Tags of endpoint.
  * `key` - (Required, String) Tag key.
  * `value` - (Optional, String) Tag value.
* `ip_address_type` - (Optional, String) IP address type. Valid values are `Ipv4` and `Ipv6`. Default is `Ipv4`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `cdc_id` - CDC instance ID.
* `create_time` - Create Time.
* `end_point_owner` - APPID.
* `state` - state of end point.


## Import

VPC end point can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_end_point.example vpce-ntv3vy9k
```

