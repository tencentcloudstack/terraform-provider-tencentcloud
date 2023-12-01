---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vpc_account_attributes"
sidebar_current: "docs-tencentcloud-datasource-vpc_account_attributes"
description: |-
  Use this data source to query detailed information of vpc account_attributes
---

# tencentcloud_vpc_account_attributes

Use this data source to query detailed information of vpc account_attributes

## Example Usage

```hcl
data "tencentcloud_vpc_account_attributes" "account_attributes" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `account_attribute_set` - User account attribute object.
  * `attribute_name` - Attribute name.
  * `attribute_values` - Attribute values.


