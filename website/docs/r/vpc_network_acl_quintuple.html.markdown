---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_network_acl_quintuple"
sidebar_current: "docs-tencentcloud-resource-vpc_network_acl_quintuple"
description: |-
  Provides a resource to create a vpc network acl quintuple
---

# tencentcloud_vpc_network_acl_quintuple

Provides a resource to create a vpc network acl quintuple

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
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones_by_product.zones.zones.0.name
}

resource "tencentcloud_vpc_acl" "example" {
  vpc_id = tencentcloud_vpc.vpc.id
  name   = "tf-example"
  ingress = [
    "ACCEPT#192.168.1.0/24#800#TCP",
    "ACCEPT#192.168.1.0/24#800-900#TCP",
  ]
  egress = [
    "ACCEPT#192.168.1.0/24#800#TCP",
    "ACCEPT#192.168.1.0/24#800-900#TCP",
  ]
}

resource "tencentcloud_vpc_network_acl_quintuple" "network_acl_quintuple" {
  network_acl_id = tencentcloud_vpc_acl.example.id
  network_acl_quintuple_set {
    ingress {
      protocol                       = "TCP"
      description                    = "ingress desc."
      source_port                    = "80"
      source_cidr                    = "192.168.0.0/24"
      destination_port               = "8080"
      destination_cidr               = "192.168.0.0/24"
      action                         = "DROP"
      network_acl_quintuple_entry_id = "acli45-q1phngkz"
      priority                       = 1
      network_acl_direction          = "INGRESS"
    }
    egress {
      protocol                       = "TCP"
      description                    = "egress desc."
      source_port                    = "80"
      source_cidr                    = "192.168.0.0/24"
      destination_port               = "8080"
      destination_cidr               = "192.168.0.0/24"
      action                         = "DROP"
      network_acl_quintuple_entry_id = "acli45-q1phngkz"
      priority                       = 1
      network_acl_direction          = "EGRESS"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `network_acl_id` - (Required, String) Network ACL instance ID. For example:acl-12345678.
* `network_acl_quintuple_set` - (Required, List) Network quintuple ACL rule set.

The `egress` object of `network_acl_quintuple_set` supports the following:

* `action` - (Optional, String) Action, ACCEPT or DROP.
* `create_time` - (Optional, String) Creation time, used as an output parameter of DescribeNetworkAclQuintupleEntries.
* `description` - (Optional, String) Description.
* `destination_cidr` - (Optional, String) Destination CIDR.
* `destination_port` - (Optional, String) Destination port (all, single port, range). When Protocol is ALL or ICMP, Port cannot be specified.
* `network_acl_direction` - (Optional, String) Direction, INGRESS or EGRESS, is used as an output parameter of DescribeNetworkAclQuintupleEntries.
* `network_acl_quintuple_entry_id` - (Optional, String) Unique ID of a network ACL entry.
* `priority` - (Optional, Int) Priority, starting from 1.
* `protocol` - (Optional, String) Protocol, value: TCP,UDP, ICMP, ALL.
* `source_cidr` - (Optional, String) Source CIDR.
* `source_port` - (Optional, String) Source port (all, single port, range). When Protocol is ALL or ICMP, Port cannot be specified.

The `ingress` object of `network_acl_quintuple_set` supports the following:

* `action` - (Optional, String) Action, ACCEPT or DROP.
* `create_time` - (Optional, String) Creation time, used as an output parameter of DescribeNetworkAclQuintupleEntries.
* `description` - (Optional, String) Description.
* `destination_cidr` - (Optional, String) Destination CIDR.
* `destination_port` - (Optional, String) Destination port (all, single port, range). When Protocol is ALL or ICMP, Port cannot be specified.
* `network_acl_direction` - (Optional, String) Direction, INGRESS or EGRESS, is used as an output parameter of DescribeNetworkAclQuintupleEntries.
* `network_acl_quintuple_entry_id` - (Optional, String) Unique ID of a network ACL entry.
* `priority` - (Optional, Int) Priority, starting from 1.
* `protocol` - (Optional, String) Protocol, value: TCP,UDP, ICMP, ALL.
* `source_cidr` - (Optional, String) Source CIDR.
* `source_port` - (Optional, String) source port (all, single port, range). When the protocol is ALL or ICMP, the port cannot be specified.

The `network_acl_quintuple_set` object supports the following:

* `egress` - (Optional, List) Network ACL quintuple outbound rule.
* `ingress` - (Optional, List) Network ACL quintuple inbound rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateNetworkAclQuintupleEntries](https://www.tencentcloud.com/document/api/215/49661)
* Read: [DescribeNetworkAclQuintupleEntries](https://www.tencentcloud.com/document/api/215/49659)
* Update: [ModifyNetworkAclQuintupleEntries](https://www.tencentcloud.com/document/api/215/49658)
* Delete: [DeleteNetworkAclQuintupleEntries](https://www.tencentcloud.com/document/api/215/49660)

## Import

vpc network_acl_quintuple can be imported using the id, e.g.

```
terraform import tencentcloud_vpc_network_acl_quintuple.network_acl_quintuple network_acl_quintuple_id
```

