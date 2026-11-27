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
  instance_id = "tdsql3-ibpyogg7tg"
  user_name   = "tf-example"
  host        = "10.0.0.1"

  global_privileges = [
    "INSERT",
    "SELECT",
    "ALTER",
    "CREATE",
    "DELETE"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `host` - (Required, String) The access host allowed for the user. The username and host uniquely identify an account.
* `instance_id` - (Required, String) Instance ID, such as tdsql3-5baee8df.
* `user_name` - (Required, String) Login username.
* `database_privileges` - (Optional, Set) Database level privileges.
* `global_privileges` - (Optional, Set) Global privileges.
* `table_privileges` - (Optional, Set) Table level privileges.

The `database_privileges` object supports the following:

* `database` - (Required, String) Database name.
* `privileges` - (Required, Set) Privilege list.

The `table_privileges` object supports the following:

* `database` - (Required, String) Database name.
* `privileges` - (Required, Set) Privilege list.
* `table` - (Required, String) Table name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the user privilege config, in the format `instance_id#user_name#host`.

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

TDSQL Boundless user privilege config can be imported using the composite id instance_id#user_name#host, e.g.

```
terraform import tencentcloud_tdmysql_user_privilege_config.example tdsql3-ibpyogg7tg#tf-example#10.0.0.1
```

