---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_eni"
sidebar_current: "docs-tencentcloud-resource-eni"
description: |-
  Provides a resource to create an ENI.
---

# tencentcloud_eni

Provides a resource to create an ENI.

## Example Usage

```hcl
data "tencentcloud_availability_zones_by_product" "zones" {
  product = "vpc"
}

resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "subnet" {
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
  name              = "subnet-example"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "10.0.0.0/16"
  is_multicast      = false
}

resource "tencentcloud_security_group" "example1" {
  name        = "tf-example-sg1"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

resource "tencentcloud_security_group" "example2" {
  name        = "tf-example-sg2"
  description = "sg desc."
  project_id  = 0

  tags = {
    "example" = "test"
  }
}

resource "tencentcloud_eni" "example" {
  name        = "tf-example-eni"
  vpc_id      = tencentcloud_vpc.vpc.id
  subnet_id   = tencentcloud_subnet.subnet.id
  description = "eni desc."
  ipv4_count  = 1
  security_groups = [
    tencentcloud_security_group.example1.id,
    tencentcloud_security_group.example2.id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the ENI, maximum length 60.
* `subnet_id` - (Required, String, ForceNew) ID of the subnet within this vpc.
* `vpc_id` - (Required, String, ForceNew) ID of the vpc.
* `description` - (Optional, String) Description of the ENI, maximum length 60.
* `ipv4_count` - (Optional, Int) The number of intranet IPv4s. When it is greater than 1, there is only one primary intranet IP. The others are auxiliary intranet IPs, which conflict with `ipv4s`.
* `ipv4s` - (Optional, Set) Applying for intranet IPv4s collection, conflict with `ipv4_count`. When there are multiple ipv4s, can only be one primary IP, and the maximum length of the array is 30. Each element contains the following attributes:
* `security_groups` - (Optional, Set: [`String`]) A set of security group IDs.
* `tags` - (Optional, Map) Tags of the ENI.

The `ipv4s` object supports the following:

* `ip` - (Required, String) Intranet IP.
* `primary` - (Required, Bool) Indicates whether the IP is primary.
* `description` - (Optional, String) Description of the IP, maximum length 25.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the ENI.
* `ipv4_info` - An information list of IPv4s. Each element contains the following attributes:
  * `description` - Description of the IP.
  * `ip` - Intranet IP.
  * `primary` - Indicates whether the IP is primary.
* `mac` - MAC address.
* `primary` - Indicates whether the IP is primary.
* `state` - State of the ENI.


## Import

ENI can be imported using the id, e.g.

```
  $ terraform import tencentcloud_eni.foo eni-qka182br
```

