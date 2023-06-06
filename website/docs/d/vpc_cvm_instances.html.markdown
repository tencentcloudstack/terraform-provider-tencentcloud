---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_cvm_instances"
sidebar_current: "docs-tencentcloud-datasource-vpc_cvm_instances"
description: |-
  Use this data source to query detailed information of vpc cvm_instances
---

# tencentcloud_vpc_cvm_instances

Use this data source to query detailed information of vpc cvm_instances

## Example Usage

```hcl
data "tencentcloud_vpc_cvm_instances" "cvm_instances" {
  filters {
    name   = "vpc-id"
    values = ["vpc-lh4nqig9"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Required, List) Filter condition. `RouteTableIds` and `Filters` cannot be specified at the same time. vpc-id - String - (Filter condition) VPC instance ID, such as `vpc-f49l6u0z`;instance-type - String - (Filter condition) CVM instance ID;instance-name - String - (Filter condition) CVM name.
* `result_output_file` - (Optional, String) Used to save results.

The `filters` object supports the following:

* `name` - (Required, String) The attribute name. If more than one Filter exists, the logical relation between these Filters is `AND`.
* `values` - (Required, Set) Attribute value. If multiple values exist in one filter, the logical relationship between these values is `OR`. For a `bool` parameter, the valid values include `TRUE` and `FALSE`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `instance_set` - List of CVM instances.
  * `cpu` - Number of CPU cores in an instance (in core).
  * `created_time` - The creation time.
  * `eni_ip_limit` - Private IP quoata for instance ENIs (including primary ENIs).
  * `eni_limit` - Instance ENI quota (including primary ENIs).
  * `instance_eni_count` - The number of ENIs (including primary ENIs) bound to a instance.
  * `instance_id` - CVM instance ID.
  * `instance_name` - CVM Name.
  * `instance_state` - CVM status.
  * `instance_type` - Instance type.
  * `memory` - Instance's memory capacity. Unit: GB.
  * `subnet_id` - Subnet instance ID.
  * `vpc_id` - VPC instance ID.


