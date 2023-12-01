---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_private_ip_addresses"
sidebar_current: "docs-tencentcloud-datasource-vpc_private_ip_addresses"
description: |-
  Use this data source to query detailed information of vpc private_ip_addresses
---

# tencentcloud_vpc_private_ip_addresses

Use this data source to query detailed information of vpc private_ip_addresses

## Example Usage

```hcl
data "tencentcloud_vpc_private_ip_addresses" "private_ip_addresses" {
  vpc_id               = "vpc-l0dw94uh"
  private_ip_addresses = ["10.0.0.1"]
}
```

## Argument Reference

The following arguments are supported:

* `private_ip_addresses` - (Required, Set: [`String`]) The private `IP` address list. Each request supports a maximum of `10` batch querying.
* `vpc_id` - (Required, String) The `ID` of the `VPC`, such as `vpc-f49l6u0z`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `vpc_private_ip_address_set` - The list of private `IP` address information.
  * `cidr_block` - The `CIDR` belonging to the subnet.
  * `created_time` - `IP` application time.
  * `private_ip_address_type` - Private `IP` type.
  * `private_ip_address` - `VPC` private `IP`.


