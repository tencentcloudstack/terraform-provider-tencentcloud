---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_acls"
sidebar_current: "docs-tencentcloud-datasource-vpc_acls"
description: |-
  Use this data source to query VPC Network ACL information.
---

# tencentcloud_vpc_acls

Use this data source to query VPC Network ACL information.

## Example Usage

```hcl
data "tencentcloud_vpc_instances" "foo" {
}

data "tencentcloud_vpc_acls" "foo" {
  vpc_id = data.tencentcloud_vpc_instances.foo.instance_list.0.vpc_id
}

data "tencentcloud_vpc_acls" "foo" {
  name = "test_acl"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, String) ID of the network ACL instance.
* `name` - (Optional, String) Name of the network ACL.
* `result_output_file` - (Optional, String) Used to save results.
* `vpc_id` - (Optional, String) ID of the VPC instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `acl_list` - The information list of the VPC. Each element contains the following attributes:
  * `create_time` - Creation time.
  * `egress` - Outbound rules of the network ACL.
    * `cidr_block` - An IP address network or segment.
    * `description` - Rule description.
    * `policy` - Rule policy of Network ACL.
    * `port` - Range of the port.
    * `protocol` - Type of IP protocol.
  * `id` - ID of the network ACL instance.
  * `ingress` - Inbound rules of the network ACL.
    * `cidr_block` - An IP address network or segment.
    * `description` - Rule description.
    * `policy` - Rule policy of Network ACL.
    * `port` - Range of the port.
    * `protocol` - Type of IP protocol.
  * `name` - Name of the network ACL.
  * `subnets` - Subnets associated with the network ACL.
    * `cidr_block` - The IPv4 CIDR of the subnet.
    * `subnet_id` - Subnet instance ID.
    * `subnet_name` - Subnet name.
    * `tags` - Tags of the subnet.
    * `vpc_id` - ID of the VPC instance.
  * `vpc_id` - ID of the VPC instance.


