---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_route_table_input_policies"
sidebar_current: "docs-tencentcloud-resource-ccn_route_table_input_policies"
description: |-
  Provides a resource to create a CCN Route table input policies.
---

# tencentcloud_ccn_route_table_input_policies

Provides a resource to create a CCN Route table input policies.

~> **NOTE:** Use this resource to manage all input policies under the routing table of CCN instances.

## Example Usage

```hcl
variable "region" {
  default = "ap-guangzhou"
}

variable "availability_zone" {
  default = "ap-guangzhou-4"
}

# create vpc
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc"
  cidr_block = "172.16.0.0/16"
}

# create subnet
resource "tencentcloud_subnet" "subnet" {
  availability_zone = var.availability_zone
  name              = "subnet"
  vpc_id            = tencentcloud_vpc.vpc.id
  cidr_block        = "172.16.0.0/24"
  is_multicast      = false
}

# create ccn
resource "tencentcloud_ccn" "example" {
  name                 = "tf-example"
  description          = "desc."
  qos                  = "AG"
  charge_type          = "PREPAID"
  bandwidth_limit_type = "INTER_REGION_LIMIT"
  tags = {
    createBy = "terraform"
  }
}

# create ccn route table
resource "tencentcloud_ccn_route_table" "example" {
  ccn_id      = tencentcloud_ccn.example.id
  name        = "tf-example"
  description = "desc."
}

# attachment instance
resource "tencentcloud_ccn_attachment" "attachment" {
  ccn_id          = tencentcloud_ccn.example.id
  instance_id     = tencentcloud_vpc.vpc.id
  instance_type   = "VPC"
  instance_region = var.region
  route_table_id  = tencentcloud_ccn_route_table.example.id
}

# create route table input policy
resource "tencentcloud_ccn_route_table_input_policies" "example" {
  ccn_id         = tencentcloud_ccn.example.id
  route_table_id = tencentcloud_ccn_route_table.example.id
  policies {
    action      = "accept"
    description = "desc."
    route_conditions {
      name          = "instance-region"
      values        = [var.region]
      match_pattern = 1
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String, ForceNew) CCN Instance ID.
* `route_table_id` - (Required, String, ForceNew) CCN Route table ID.
* `policies` - (Optional, List) Routing reception strategy.

The `policies` object supports the following:

* `action` - (Required, String) Routing behavior, `accept` allows, `drop` rejects.
* `description` - (Required, String) Policy description.
* `route_conditions` - (Required, List) Routing conditions.

The `route_conditions` object of `policies` supports the following:

* `match_pattern` - (Required, Int) Matching mode, `1` precise matching, `0` fuzzy matching.
* `name` - (Required, String) Condition type. Example value: `instance-type`, `instance-region`, `instance-id`, `cidr-block`.
* `values` - (Required, List) List of conditional values. Example value:
 `instance-type`: `VPC`, `VPNGW`, `DIRECTCONNECT`
 `instance-region`: `ap-guangzhou`
 `instance-id`: `vpc-axrsmmrv`, `dcg-oxad32f7`, `vpngw-33p5vnwd`
 `cidr-block`: `172.0.0.0/8`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Ccn instance can be imported, e.g.

```
$ terraform import tencentcloud_ccn_route_table_input_policies.example ccn-gr7nynbd#ccnrtb-jpf7bzn3
```

