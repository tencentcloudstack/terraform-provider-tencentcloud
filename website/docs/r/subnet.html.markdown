---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_subnet"
sidebar_current: "docs-tencentcloud-resource-subnet"
description: |-
  Provide a resource to create a VPC subnet.
---

# tencentcloud_subnet

Provide a resource to create a VPC subnet.

~> **NOTE:** In accordance with VPC business requirements, the default value for `is_multicast` has been updated to `false`(previously `true`) in version `v1.82.93` of the provider. If you wish to utilize this feature, you must first contact the VPC product team to have your account added to the whitelist, and then set the `is_multicast` field to `true`.

## Example Usage

### Create a normal VPC subnet

```hcl
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name         = "vpc-example"
  cidr_block   = "10.0.0.0/16"
  is_multicast = false
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  is_multicast      = false
}
```

### Create a CDC instance VPC subnet

```hcl
data "tencentcloud_availability_zones" "zones" {}

resource "tencentcloud_vpc" "vpc" {
  name         = "vpc-example"
  cidr_block   = "10.0.0.0/16"
  is_multicast = false
}

resource "tencentcloud_subnet" "subnet" {
  vpc_id            = tencentcloud_vpc.vpc.id
  name              = "subnet-example"
  cidr_block        = "10.0.0.0/16"
  cdc_id            = "cluster-lchwgxhs"
  availability_zone = data.tencentcloud_availability_zones.zones.zones.0.name
  is_multicast      = false
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required, String, ForceNew) The availability zone within which the subnet should be created.
* `cidr_block` - (Required, String, ForceNew) A network address block of the subnet.
* `name` - (Required, String) The name of subnet to be created.
* `vpc_id` - (Required, String, ForceNew) ID of the VPC to be associated.
* `cdc_id` - (Optional, String, ForceNew) ID of CDC instance.
* `is_multicast` - (Optional, Bool) Indicates whether multicast is enabled. The default value is `false`. We recommend disabling these features if they are not applicable to your environment.
* `route_table_id` - (Optional, String) ID of a routing table to which the subnet should be associated.
* `tags` - (Optional, Map) Tags of the subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `available_ip_count` - The number of available IPs.
* `create_time` - Creation time of subnet resource.
* `is_default` - Indicates whether it is the default VPC for this region.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateSubnet](https://cloud.tencent.com/document/api/215/15782) - This API is used to create a subnet.
* Create: [DescribeRouteTables](https://cloud.tencent.com/document/api/215/15763) - This API is used to query route tables.
* Create: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Create: [ModifySubnetAttribute](https://cloud.tencent.com/document/api/215/15781) - This API (ModifySubnetAttribute) is used to modify subnet attributes.
* Create: [ReplaceRouteTableAssociation](https://cloud.tencent.com/document/api/215/15767) - This API (ReplaceRouteTableAssociation) is used to modify the route table associated with a subnet.
* Read: [DescribeResourceTagsByResourceIds](https://cloud.tencent.com/document/api/651/35319) - This API is used to query the tag key-value pairs associated with an existing resource.
* Read: [DescribeSubnets](https://cloud.tencent.com/document/api/215/15784) - This API (DescribeSubnets) is used to query the list of subnets.
* Update: [DescribeRouteTables](https://cloud.tencent.com/document/api/215/15763) - This API is used to query route tables.
* Update: [ModifyResourceTags](https://cloud.tencent.com/document/api/651/35322) - This API is used to modify all tags associated with a resource.
* Update: [ModifySubnetAttribute](https://cloud.tencent.com/document/api/215/15781) - This API (ModifySubnetAttribute) is used to modify subnet attributes.
* Update: [ReplaceRouteTableAssociation](https://cloud.tencent.com/document/api/215/15767) - This API (ReplaceRouteTableAssociation) is used to modify the route table associated with a subnet.
* Delete: [DeleteSubnet](https://cloud.tencent.com/document/api/215/15783) - This API is used to delete a subnet.

## Import

Vpc subnet instance can be imported, e.g.

```
$ terraform import tencentcloud_subnet.subnet subnet-b8j03v0c
```

