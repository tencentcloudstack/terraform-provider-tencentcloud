---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc"
sidebar_current: "docs-tencentcloud-resource-vpc"
description: |-
  Provide a resource to create a VPC.
---

# tencentcloud_vpc

Provide a resource to create a VPC.

## Example Usage

### Create a basic VPC

```hcl
resource "tencentcloud_vpc" "vpc" {
  name         = "tf-example"
  cidr_block   = "10.0.0.0/16"
  dns_servers  = ["119.29.29.29", "8.8.8.8"]
  is_multicast = false

  tags = {
    "test" = "test"
  }
}
```

### Using Assistant CIDR

```hcl
resource "tencentcloud_vpc" "vpc" {
  name            = "tf-example"
  cidr_block      = "10.0.0.0/16"
  is_multicast    = false
  assistant_cidrs = ["172.16.0.0/24"]

  tags = {
    "test" = "test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Required, String, ForceNew) A network address block which should be a subnet of the three internal network segments (10.0.0.0/16, 172.16.0.0/12 and 192.168.0.0/16).
* `name` - (Required, String) The name of the VPC.
* `assistant_cidrs` - (Optional, List: [`String`]) List of Assistant CIDR, NOTE: Only `NORMAL` typed CIDRs included, check the Docker CIDR by readonly `assistant_docker_cidrs`.
* `dns_servers` - (Optional, Set: [`String`]) The DNS server list of the VPC. And you can specify 0 to 5 servers to this list.
* `is_multicast` - (Optional, Bool) Indicates whether VPC multicast is enabled. The default value is 'true'.
* `tags` - (Optional, Map) Tags of the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of VPC.
* `default_route_table_id` - Default route table id, which created automatically after VPC create.
* `docker_assistant_cidrs` - List of Docker Assistant CIDR.
* `is_default` - Indicates whether it is the default VPC for this region.


## Import

Vpc instance can be imported, e.g.

```
$ terraform import tencentcloud_vpc.test vpc-id
```

