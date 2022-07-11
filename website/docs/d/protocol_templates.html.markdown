---
subcategory: "Virtual Private Cloud(VPC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_protocol_templates"
sidebar_current: "docs-tencentcloud-datasource-protocol_templates"
description: |-
  Use this data source to query detailed information of protocol templates.
---

# tencentcloud_protocol_templates

Use this data source to query detailed information of protocol templates.

## Example Usage

```hcl
data "tencentcloud_protocol_templates" "name" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, String) ID of the protocol template to query.
* `name` - (Optional, String) Name of the protocol template to query.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `template_list` - Information list of the dedicated protocol templates.
  * `id` - ID of the protocol template.
  * `name` - Name of protocol template.
  * `protocols` - Set of the protocols.


