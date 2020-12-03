---
subcategory: "VPC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_address_template_groups"
sidebar_current: "docs-tencentcloud-datasource-address_template_groups"
description: |-
  Use this data source to query detailed information of address template groups.
---

# tencentcloud_address_template_groups

Use this data source to query detailed information of address template groups.

## Example Usage

```hcl
data "tencentcloud_address_template" "name" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) Id of the address template group to query.
* `name` - (Optional) Name of the address template group to query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `group_list` - Information list of the dedicated address template groups.
  * `id` - Id of the address template group.
  * `name` - Name of address template group.
  * `template_ids` - ID set of the address template.


