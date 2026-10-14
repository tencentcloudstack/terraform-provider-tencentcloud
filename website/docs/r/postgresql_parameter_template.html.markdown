---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_parameter_template"
sidebar_current: "docs-tencentcloud-resource-postgresql_parameter_template"
description: |-
  Provides a resource to create a PostgreSQL parameter template
---

# tencentcloud_postgresql_parameter_template

Provides a resource to create a PostgreSQL parameter template

## Example Usage

```hcl
resource "tencentcloud_postgresql_parameter_template" "example" {
  template_name        = "tf-example"
  db_major_version     = "18"
  db_engine            = "postgresql"
  template_description = "remark."

  modify_param_entry_set {
    name           = "timezone"
    expected_value = "PRC"
  }

  modify_param_entry_set {
    name           = "lock_timeout"
    expected_value = "60"
  }

  modify_param_entry_set {
    name           = "event_triggers"
    expected_value = "on"
  }
}
```



```hcl
resource "tencentcloud_postgresql_parameter_template" "example" {
  template_name        = "tf-example"
  db_major_version     = "18"
  db_engine            = "postgresql"
  template_description = "remark."

  modify_param_entry_set {
    name           = "timezone"
    expected_value = "PRC"
  }

  modify_param_entry_set {
    name           = "event_triggers"
    expected_value = "on"
  }

  delete_param_set = ["lock_timeout"]
}
```

## Argument Reference

The following arguments are supported:

* `db_engine` - (Required, String, ForceNew) Database engine, such as postgresql, mssql_compatible.
* `db_major_version` - (Required, String, ForceNew) The major database version number, such as 11, 12, 13.
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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateParameterTemplate](https://www.tencentcloud.com/document/product/409/52655) - This API is used to create a parameter template.
* Create: [ModifyParameterTemplate](https://www.tencentcloud.com/document/product/409/52650) - This API is used to modify the configurations, such as parameter template name and description. It can also be used to manage the parameter list in the parameter template.
* Read: [DescribeParameterTemplateAttributes](https://www.tencentcloud.com/document/product/409/52652) - This API is used to query the details of a parameter template, including basic information and parameter information.
* Update: [ModifyParameterTemplate](https://www.tencentcloud.com/document/product/409/52650) - This API is used to modify the configurations, such as parameter template name and description. It can also be used to manage the parameter list in the parameter template.
* Delete: [DeleteParameterTemplate](https://www.tencentcloud.com/document/product/409/52654) - This API is used to delete a parameter template.

## Import

PostgreSQL parameter template can be imported using the id, e.g.

Notice: `delete_param_set` do not support import.

```
terraform import tencentcloud_postgresql_parameter_template.example 81ec47ed-0e4e-5af2-a648-2072fe63f225
```

