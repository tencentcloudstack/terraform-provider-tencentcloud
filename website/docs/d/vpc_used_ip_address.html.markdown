---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_used_ip_address"
sidebar_current: "docs-tencentcloud-datasource-vpc_used_ip_address"
description: |-
  Use this data source to query detailed information of vpc used_ip_address
---

# tencentcloud_vpc_used_ip_address

Use this data source to query detailed information of vpc used_ip_address

## Example Usage

```hcl
data "tencentcloud_vpc_used_ip_address" "used_ip_address" {
  vpc_id = "vpc-4owdpnwr"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, String) VPC instance ID.
* `ip_addresses` - (Optional, Set: [`String`]) IPs to query.
* `result_output_file` - (Optional, String) Used to save results.
* `subnet_id` - (Optional, String) Subnet instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ip_address_states` - Information of resources bound with the queried IPs Note: This parameter may return null, indicating that no valid values can be obtained.
  * `ip_address` - IP address.
  * `resource_id` - Resource ID.
  * `resource_type` - Resource type.
  * `subnet_id` - Subnet instance ID.
  * `vpc_id` - VPC instance ID.


