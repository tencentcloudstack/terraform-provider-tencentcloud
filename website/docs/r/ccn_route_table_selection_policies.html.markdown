---
subcategory: "Cloud Connect Network(CCN)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ccn_route_table_selection_policies"
sidebar_current: "docs-tencentcloud-resource-ccn_route_table_selection_policies"
description: |-
  Provides a resource to create a CCN Route table selection policies.
---

# tencentcloud_ccn_route_table_selection_policies

Provides a resource to create a CCN Route table selection policies.

~> **NOTE:** Use this resource to manage all selection policies under the routing table of CCN instances.

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

# create route table selection policy
resource "tencentcloud_ccn_route_table_selection_policies" "example" {
  ccn_id = tencentcloud_ccn.example.id
  selection_policies {
    instance_type     = "VPC"
    instance_id       = tencentcloud_vpc.vpc.id
    source_cidr_block = "192.168.100.0/24"
    route_table_id    = tencentcloud_ccn_route_table.example.id
    description       = "desc."
  }
}
```

## Argument Reference

The following arguments are supported:

* `ccn_id` - (Required, String, ForceNew) CCN Instance ID.
* `selection_policies` - (Required, List) Select strategy information set.

The `selection_policies` object supports the following:

* `description` - (Required, String) description.
* `instance_id` - (Required, String) Instance ID.
* `instance_type` - (Required, String) Instance Type: Private Network: VPC, Dedicated Gateway: DIRECTCONNECT, Blackstone Private Network: BMVPC, EDGE Device: EDGE, EDGE Tunnel: EDGE_TUNNEL, EDGE Gateway: EDGE_VPNGW, VPN Gateway: VPNGW.
* `route_table_id` - (Required, String) route table ID.
* `source_cidr_block` - (Required, String) Source CIDR.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Ccn instance can be imported, e.g.

```
$ terraform import tencentcloud_ccn_route_table_selection_policies.example ccn-gr7nynbd
```

