---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_param_template"
sidebar_current: "docs-tencentcloud-resource-mysql_param_template"
description: |-
  Provides a resource to create a mysql param template
---

# tencentcloud_mysql_param_template

Provides a resource to create a mysql param template

## Example Usage

```hcl
resource "tencentcloud_mysql_param_template" "param_template" {
  name           = "terraform-test"
  description    = "terraform-test"
  engine_version = "8.0"
  template_type  = "HIGH_STABILITY"
  engine_type    = "InnoDB"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of parameter template.
* `description` - (Optional, String) The description of parameter template.
* `engine_type` - (Optional, String) The engine type of instance, optional value is InnoDB or RocksDB, default to InnoDB.
* `engine_version` - (Optional, String) The version of MySQL.
* `template_type` - (Optional, String) The default type of parameter template, supported value is HIGH_STABILITY or HIGH_PERFORMANCE.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `param_list` - parameter list.
  * `current_value` - The value of parameter.
  * `name` - The name of parameter.
* `template_id` - The ID of source parameter template.


## Import

mysql param template can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_param_template.param_template template_id
```

