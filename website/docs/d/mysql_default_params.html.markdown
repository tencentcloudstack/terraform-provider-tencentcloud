---
subcategory: "MySQL"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_default_params"
sidebar_current: "docs-tencentcloud-datasource-mysql_default_params"
description: |-
  Provide a datasource to query default mysql parameters.
---

# tencentcloud_mysql_default_params

Provide a datasource to query default mysql parameters.

## Example Usage

```hcl
resource "tencentcloud_mysql_default_params" "mysql_57" {
  db_version = "5.7"
}
```

## Argument Reference

The following arguments are supported:

* `db_version` - (Optional, String) MySQL database version.
* `result_output_file` - (Optional, String) Used for save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `param_list` - List of param detail.
  * `current_value` - Param current value.
  * `default` - Param default value.
  * `description` - Param description.
  * `enum_value` - Params available values if type of param is enum.
  * `max` - Param maximum value if type of param is integer.
  * `min` - Param minimum value if type of param is integer.
  * `name` - Param key name.
  * `need_reboot` - Indicates weather the database instance reboot if param modified.
  * `param_type` - Type of param.


