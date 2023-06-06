---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_classic_link_instances"
sidebar_current: "docs-tencentcloud-datasource-vpc_classic_link_instances"
description: |-
  Use this data source to query detailed information of vpc classic_link_instances
---

# tencentcloud_vpc_classic_link_instances

Use this data source to query detailed information of vpc classic_link_instances

## Example Usage

```hcl
data "tencentcloud_vpc_classic_link_instances" "classic_link_instances" {
  filters {
    name   = "vpc-id"
    values = ["vpc-lh4nqig9"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions.`vpc-id` - String - (Filter condition) The VPC instance ID. `vm-ip` - String - (Filter condition) The IP address of the CVM on the basic network.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) The attribute name. If more than one Filter exists, the logical relation between these Filters is `AND`.
* `values` - (Required, Set) The attribute value. If there are multiple Values for one Filter, the logical relation between these Values under the same Filter is `OR`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `classic_link_instance_set` - Classiclink instance.
  * `instance_id` - The unique ID of the CVM instance.
  * `vpc_id` - VPC instance ID.


