---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_subnet_resource_dashboard"
sidebar_current: "docs-tencentcloud-datasource-vpc_subnet_resource_dashboard"
description: |-
  Use this data source to query detailed information of vpc subnet_resource_dashboard
---

# tencentcloud_vpc_subnet_resource_dashboard

Use this data source to query detailed information of vpc subnet_resource_dashboard

## Example Usage

```hcl
data "tencentcloud_vpc_subnet_resource_dashboard" "subnet_resource_dashboard" {
  subnet_ids = ["subnet-i9tpf6hq"]
}
```

## Argument Reference

The following arguments are supported:

* `subnet_ids` - (Required, Set: [`String`]) Subnet instance ID, such as `subnet-f1xjkw1b`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `resource_statistics_set` - Information of resources returned.
  * `ip` - The total number of used IP addresses.
  * `resource_statistics_item_set` - Information of associated resources.
    * `resource_count` - Number of resources.
    * `resource_name` - Resource name.
    * `resource_type` - Resource type, such as CVM, ENI.
  * `subnet_id` - Subnet instance ID, such as `subnet-bthucmmy`.
  * `vpc_id` - VPC instance ID, such as vpc-f1xjkw1b.


