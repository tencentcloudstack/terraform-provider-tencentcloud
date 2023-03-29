---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_parameter_template"
sidebar_current: "docs-tencentcloud-resource-postgresql_parameter_template"
description: |-
  Provides a resource to create a postgresql parameter_template
---

# tencentcloud_postgresql_parameter_template

Provides a resource to create a postgresql parameter_template

## Example Usage

```hcl
resource "tencentcloud_postgresql_parameter_template" "parameter_template" {
  template_name        = "your_temp_name"
  db_major_version     = "13"
  db_engine            = "postgresql"
  template_description = "For_tf_test"

  modify_param_entry_set {
    name           = "timezone"
    expected_value = "UTC"
  }
  modify_param_entry_set {
    name           = "lock_timeout"
    expected_value = "123"
  }

  delete_param_set = ["lc_time"]
}
```

## Argument Reference

The following arguments are supported:

* `db_engine` - (Required, String) Database engine, such as postgresql, mssql_compatible.
* `db_major_version` - (Required, String) The major database version number, such as 11, 12, 13.
* `template_name` - (Required, String) Template name, which can contain 1-60 letters, digits, and symbols (-_./()+=:@).
* `delete_param_set` - (Optional, Set: [`String`]) The set of parameters that need to be deleted.
* `modify_param_entry_set` - (Optional, Set) The set of parameters that need to be modified or added. Note: the same parameter cannot appear in the set of modifying and adding and deleting at the same time.
* `template_description` - (Optional, String) Parameter template description, which can contain 1-60 letters, digits, and symbols (-_./()+=:@).

The `modify_param_entry_set` object supports the following:

* `expected_value` - (Required, String) Modify the parameter value. The input parameters are passed in the form of strings, for example: decimal `0.1`, integer `1000`, enumeration `replica`.
* `name` - (Required, String) The parameter name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

postgresql parameter_template can be imported using the id, e.g.

Notice: `modify_param_entry_set` and `delete_param_set` do not support import.

```
terraform import tencentcloud_postgresql_parameter_template.parameter_template parameter_template_id
```

