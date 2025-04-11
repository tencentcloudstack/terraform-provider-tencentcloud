---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_acl"
sidebar_current: "docs-tencentcloud-resource-vpc_acl"
description: |-
  Provide a resource to create a VPC ACL instance.
---

# tencentcloud_vpc_acl

Provide a resource to create a VPC ACL instance.

## Example Usage

```hcl
resource "tencentcloud_vpc" "vpc" {
  name       = "vpc-example"
  cidr_block = "10.0.0.0/16"
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Name of the network ACL.
* `vpc_id` - (Required, String) ID of the VPC instance.
* `egress` - (Optional, List: [`String`]) Egress rules. A rule must match the following format: [action]#[cidr_ip]#[port]#[protocol]. The available value of 'action' is `ACCEPT` and `DROP`. The 'cidr_ip' must be an IP address network or segment. The 'port' valid format is `80`, `80-90` or `ALL`. The available value of 'protocol' is `TCP`, `UDP`, `ICMP` and `ALL`. When 'protocol' is `ICMP` or `ALL`, the 'port' must be `ALL`.
* `ingress` - (Optional, List: [`String`]) Ingress rules. A rule must match the following format: [action]#[cidr_ip]#[port]#[protocol]. The available value of 'action' is `ACCEPT` and `DROP`. The 'cidr_ip' must be an IP address network or segment. The 'port' valid format is `80`, `80-90` or `ALL`. The available value of 'protocol' is `TCP`, `UDP`, `ICMP` and `ALL`. When 'protocol' is `ICMP` or `ALL`, the 'port' must be `ALL`.
* `tags` - (Optional, Map) Tags of the vpc acl.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of ACL.


## Import

Vpc ACL can be imported, e.g.

```
$ terraform import tencentcloud_vpc_acl.example acl-mijp2oyu
```

