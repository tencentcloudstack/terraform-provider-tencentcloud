---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_subnet"
sidebar_current: "docs-tencentcloud-datasource-subnet"
description: |-
  Provides details about a specific VPC subnet.
---

# tencentcloud_subnet

Provides details about a specific VPC subnet.

This resource can prove useful when a module accepts a subnet id as an input variable and needs to, for example, determine the id of the VPC that the subnet belongs to.

~> **NOTE:** It has been deprecated and replaced by tencentcloud_vpc_subnets.

## Example Usage

### Query method 1

```hcl
data "tencentcloud_subnet" "subnet" {
  vpc_id    = "vpc-ha5l97e3"
  subnet_id = "subnet-ezgfompo"
}
```

### Query method 2

```hcl
data "tencentcloud_subnet" "subnet" {
  vpc_id    = "vpc-ha5l97e3"
  subnet_id = "subnet-ezgfompo"
  cdc_id    = "cluster-lchwgxhs"
}
```

## Argument Reference

The following arguments are supported:

* `subnet_id` - (Required, String) The ID of the Subnet.
* `vpc_id` - (Required, String) The VPC ID.
* `cdc_id` - (Optional, String) ID of CDC instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `availability_zone` - The AZ for the subnet.
* `cidr_block` - The CIDR block of the Subnet.
* `name` - The name for the Subnet.
* `route_table_id` - The Route Table ID.


