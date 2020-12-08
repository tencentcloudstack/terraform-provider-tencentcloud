---
subcategory: "VPC"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_service_templates"
sidebar_current: "docs-tencentcloud-datasource-service_templates"
description: |-
  Use this data source to query detailed information of service templates.
---

# tencentcloud_service_templates

Use this data source to query detailed information of service templates.

## Example Usage

```hcl
data "tencentcloud_service_template" "name" {
  name = "test"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) Id of the service template to query.
* `name` - (Optional) Name of the service template to query.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `template_list` - Information list of the dedicated service templates.
  * `id` - Id of the service template.
  * `name` - Name of service template.
  * `services` - Set of the services.


