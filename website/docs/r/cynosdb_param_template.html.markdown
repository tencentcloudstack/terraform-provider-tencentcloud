---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_param_template"
sidebar_current: "docs-tencentcloud-resource-cynosdb_param_template"
description: |-
  Provides a resource to create a cynosdb param_template
---

# tencentcloud_cynosdb_param_template

Provides a resource to create a cynosdb param_template

## Example Usage

```hcl
resource "tencentcloud_cynosdb_param_template" "param_template" {
  db_mode              = "SERVERLESS"
  engine_version       = "5.7"
  template_description = "terraform-template"
  template_name        = "terraform-template"

  param_list {
    current_value = "-1"
    param_name    = "optimizer_trace_offset"
  }
}
```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, String) MySQL version number.
* `template_name` - (Required, String) Template Name.
* `db_mode` - (Optional, String) Database type, optional values: NORMAL (default), SERVERLESS.
* `param_list` - (Optional, Set) parameter list.
* `template_description` - (Optional, String) Template Description.
* `template_id` - (Optional, Int) Optional parameter, template ID to be copied.

The `param_list` object supports the following:

* `current_value` - (Optional, String) Current value.
* `param_name` - (Optional, String) Parameter Name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



