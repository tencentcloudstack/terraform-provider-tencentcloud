---
subcategory: "Virtual Private Cloud(VPC)"
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
data "tencentcloud_address_templates" "name" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, String) ID of the address template to query.
* `name` - (Optional, String) Name of the address template to query.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `template_list` - Information list of the dedicated address templates.
  * `addresses` - Set of the addresses.
  * `id` - ID of the address template.
  * `name` - Name of address template.


