---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_protocol_template_groups"
sidebar_current: "docs-tencentcloud-datasource-protocol_template_groups"
description: |-
  Use this data source to query detailed information of protocol template groups.
---

# tencentcloud_protocol_template_groups

Use this data source to query detailed information of protocol template groups.

## Example Usage

```hcl
data "tencentcloud_protocol_template_groups" "name" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) Id of the protocol template group to query.
* `name` - (Optional) Name of the protocol template group to query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `group_list` - Information list of the dedicated protocol template groups.
  * `id` - Id of the protocol template group.
  * `name` - Name of protocol template group.
  * `template_ids` - ID set of the protocol template.


