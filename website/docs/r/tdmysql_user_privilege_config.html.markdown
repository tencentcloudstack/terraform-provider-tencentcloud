---
subcategory: "TDSQL Boundless(tdmysql)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmysql_user_privilege_config_config"
sidebar_current: "docs-tencentcloud-resource-tdmysql_user_privilege_config"
description: |-
  Provides a resource to manage user privileges of a TDSQL Boundless instance.
---

# tencentcloud_tdmysql_user_privilege_config_config

Provides a resource to manage user privileges of a TDSQL Boundless instance.

## Example Usage

```hcl
resource "tencentcloud_tdmysql_user_privilege_config_config" "example" {
  instance_id = "tdsql3-5baee8df"
  host        = "localhost"
  user_name   = "tf-example-user"
  db_name     = "*"
  privileges  = ["SELECT", "INSERT"]
}
```

## Argument Reference

The following arguments are supported:

* `db_name` - (Required, String) Database name. If it is `*`, it represents global privileges (that is, `*.*`), in which case the object_type and object parameters are ignored.
* `host` - (Required, String) The access host allowed for the user. The username and host uniquely identify an account.
* `instance_id` - (Required, String) Instance ID, such as tdsql3-5baee8df.
* `privileges` - (Required, List) Privilege list.
* `user_name` - (Required, String) Login username.
* `col_name` - (Optional, String) When object_type is table, col_name as `*` means querying the privileges of the table, and a specific column name means querying the privileges of the corresponding column.
* `object_type` - (Optional, String) Type, which can be table, view, proc, func and `*`. When db_name is a concrete database name and object_type is `*`, it represents the privileges of that database (that is, db.*), in which case the object parameter is ignored.
* `object` - (Optional, String) The name of the specific type, for example the specific table name when object_type is table. When db_name and object_type are both concrete names, object represents the specific object name, which cannot be `*` or empty.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the user privilege, in the format `instance_id#host#user_name#db_name`.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create:
  - [DescribeUserPrivileges](https://www.tencentcloud.com/document/api/1292/80120)
  - [ModifyUserPrivileges](https://www.tencentcloud.com/document/api/1292/80104)
* Read: [DescribeUserPrivileges](https://www.tencentcloud.com/document/api/1292/80120)
* Update:
  - [DescribeUserPrivileges](https://www.tencentcloud.com/document/api/1292/80120)
  - [ModifyUserPrivileges](https://www.tencentcloud.com/document/api/1292/80104)

## Import

TDSQL Boundless user privilege can be imported using the composite id instance_id#host#user_name#db_name, e.g.

```
terraform import tencentcloud_tdmysql_user_privilege_config_config.example tdsql3-5baee8df#localhost#tf-example-user#*
```

