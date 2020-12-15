---
subcategory: "VPC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_service_template_groups"
sidebar_current: "docs-tencentcloud-datasource-service_template_groups"
description: |-
  Use this data source to query detailed information of service template groups.
---

# tencentcloud_service_template_groups

Use this data source to query detailed information of service template groups.

## Example Usage

```hcl
data "tencentcloud_service_template" "name" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) Id of the service template group to query.
* `name` - (Optional) Name of the service template group to query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `group_list` - Information list of the dedicated service template groups.
  * `id` - Id of the service template group.
  * `name` - Name of service template group.
  * `template_ids` - ID set of the service template.


