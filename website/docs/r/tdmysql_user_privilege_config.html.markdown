---
subcategory: "TDSQL Boundless(tdmysql)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tdmysql_user_privilege_config"
sidebar_current: "docs-tencentcloud-resource-tdmysql_user_privilege_config"
description: |-
  Provides a resource to manage user privileges of a TDSQL Boundless instance.
---

# tencentcloud_tdmysql_user_privilege_config

Provides a resource to manage user privileges of a TDSQL Boundless instance.

## Example Usage

```hcl
resource "tencentcloud_tdmysql_user_privilege_config" "example" {
  instance_id = "tdsql3-5baee8df"
  user_name   = "tf-example-user"
  host        = "localhost"

  global_privileges = ["SELECT", "INSERT"]

  database_privileges = [
    {
      database   = "tf-example-db"
      privileges = ["SELECT"]
    },
  ]

  table_privileges = [
    {
      database   = "tf-example-db"
      table      = "tf-example-table"
      privileges = ["SELECT", "UPDATE"]
    },
  ]
}
```

## Argument Reference

The following arguments are supported:

* `database_privileges` - (Optional, List) Database level privileges.
  * `database` - (Required, String) Database name.
  * `privileges` - (Required, List) Privilege list.
* `global_privileges` - (Optional, List) Global privileges.
* `host` - (Required, String) The access host allowed for the user. The username and host uniquely identify an account.
* `instance_id` - (Required, String) Instance ID, such as tdsql3-5baee8df.
* `table_privileges` - (Optional, List) Table level privileges.
  * `database` - (Required, String) Database name.
  * `privileges` - (Required, List) Privilege list.
  * `table` - (Required, String) Table name.
* `user_name` - (Required, String) Login username.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the user privilege config, in the format `instance_id#user_name#host`.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [ModifyUserPrivileges](https://www.tencentcloud.com/document/api/1292/80104)
* Read: [DescribeUserPrivileges](https://www.tencentcloud.com/document/api/1292/80120)
* Update: [ModifyUserPrivileges](https://www.tencentcloud.com/document/api/1292/80104)

## Import

TDSQL Boundless user privilege config can be imported using the composite id instance_id#user_name#host, e.g.

```
terraform import tencentcloud_tdmysql_user_privilege_config.example tdsql3-5baee8df#tf-example-user#localhost
```
