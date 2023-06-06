---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_security_group_references"
sidebar_current: "docs-tencentcloud-datasource-vpc_security_group_references"
description: |-
  Use this data source to query detailed information of vpc security_group_references
---

# tencentcloud_vpc_security_group_references

Use this data source to query detailed information of vpc security_group_references

## Example Usage

```hcl
data "tencentcloud_vpc_security_group_references" "security_group_references" {
  security_group_ids = ["sg-edmur627"]
}
```

## Argument Reference

The following arguments are supported:

* `security_group_ids` - (Required, Set: [`String`]) A set of security group instance IDs, e.g. [sg-12345678].
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `referred_security_group_set` - Referred security groups.
  * `referred_security_group_ids` - IDs of all referred security group instances.
  * `security_group_id` - Security group instance ID.


