---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_parameter_list"
sidebar_current: "docs-tencentcloud-tencentcloud_mysql_parameter_list"
description: |-
 Use this data source to get information about a parameter group of a database instance.
---

# tencentcloud_mysql_parameter_list

Use this data source to get information about a parameter group of a database instance.

## Example Usage

```hcl
data "tencentcloud_mysql_parameter_list" "mysql" {
    mysql_id = "my-test-database "
    engine_version = "5.5" 
}

```
## Argument Reference

The following arguments are supported:

- `mysql_id` - (Optional) Instance ID. 
- `engine_version` – (Optional) The version number of the database engine to use. Supported versions include 5.5/5.6/5.7.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `parameter_name` - Parameter name.
- `parameter_type` - Parameter type.
- `description` - Parameter specification description.  
- `current_value` - Current value.
- `default_value` – Default value.
- `max` - Maximum value for the parameter.
- `min` - Minimum value for the parameter.
- `enum_value` - enumerated value.
- `need_reboot` - Indicates whether reboot is needed to enable the new parameters.