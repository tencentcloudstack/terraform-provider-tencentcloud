---
subcategory: "VPC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_address_templates"
sidebar_current: "docs-tencentcloud-datasource-address_templates"
description: |-
  Use this data source to query detailed information of address templates.
---

# tencentcloud_address_templates

Use this data source to query detailed information of address templates.

## Example Usage

```hcl
data "tencentcloud_address_template" "name" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) Id of the address template to query.
* `name` - (Optional) Name of the address template to query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `template_list` - Information list of the dedicated address templates.
  * `addresses` - Set of the addresses.
  * `id` - Id of the address template.
  * `name` - Name of address template.


