---
subcategory: "VPC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_route_table_entry"
sidebar_current: "docs-tencentcloud-resource-route_table_entry"
description: |-
  Provides a resource to create an entry of a routing table.
---

# tencentcloud_route_table_entry

Provides a resource to create an entry of a routing table.

## Example Usage

```hcl
variable "availability_zone" {
  default = "na-siliconvalley-1"
}

resource "tencentcloud_vpc" "foo" {
  name       = "ci-temp-test"
  cidr_block = "10.0.0.0/16"
}

resource "tencentcloud_subnet" "foo" {
  vpc_id            = tencentcloud_vpc.foo.id
  name              = "terraform test subnet"
  cidr_block        = "10.0.12.0/24"
  availability_zone = var.availability_zone
  route_table_id    = tencentcloud_route_table.foo.id
}

resource "tencentcloud_route_table" "foo" {
  vpc_id = tencentcloud_vpc.foo.id
  name   = "ci-temp-test-rt"
}

resource "tencentcloud_route_table_entry" "instance" {
  route_table_id         = tencentcloud_route_table.foo.id
  destination_cidr_block = "10.4.4.0/24"
  next_type              = "EIP"
  next_hub               = "0"
  description            = "ci-test-route-table-entry"
}
```

## Argument Reference

The following arguments are supported:

* `destination_cidr_block` - (Required, ForceNew) Destination address block.
* `next_hub` - (Required, ForceNew) ID of next-hop gateway. Note: when 'next_type' is EIP, GatewayId should be '0'.
* `next_type` - (Required, ForceNew) Type of next-hop, and available values include CVM, VPN, DIRECTCONNECT, PEERCONNECTION, SSLVPN, NAT, NORMAL_CVM, EIP and CCN.
* `route_table_id` - (Required, ForceNew) ID of routing table to which this entry belongs.
* `description` - (Optional, ForceNew) Description of the routing table entry.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

Route table entry can be imported using the id, e.g.

```
$ terraform import tencentcloud_route_table_entry.foo 83517.rtb-mlhpg09u
```

