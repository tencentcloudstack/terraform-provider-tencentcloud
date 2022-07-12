---
subcategory: "MySQL"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_parameter_list"
sidebar_current: "docs-tencentcloud-datasource-mysql_parameter_list"
description: |-
  Use this data source to get information about a parameter group of a database instance.
---

# tencentcloud_mysql_parameter_list

Use this data source to get information about a parameter group of a database instance.

## Example Usage

```hcl
data "tencentcloud_mysql_parameter_list" "mysql" {
  mysql_id           = "my-test-database"
  engine_version     = "5.5"
  result_output_file = "mytestpath"
}
```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Optional, String) The version number of the database engine to use. Supported versions include 5.5/5.6/5.7.
* `mysql_id` - (Optional, String) Instance ID.
* `result_output_file` - (Optional, String) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `parameter_list` - A list of parameters. Each element contains the following attributes:
  * `current_value` - Current value.
  * `default_value` - Default value.
  * `description` - Parameter specification description.
  * `enum_value` - Enumerated value.
  * `max` - Maximum value for the parameter.
  * `min` - Minimum value for the parameter.
  * `need_reboot` - Indicates whether reboot is needed to enable the new parameters.
  * `parameter_name` - Parameter name.
  * `parameter_type` - Parameter type.


