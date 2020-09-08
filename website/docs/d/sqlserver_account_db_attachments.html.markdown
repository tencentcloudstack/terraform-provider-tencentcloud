---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_account_db_attachments"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_account_db_attachments"
description: |-
  Use this data source to query the list of SQL Server account DB privileges.
---

# tencentcloud_sqlserver_account_db_attachments

Use this data source to query the list of SQL Server account DB privileges.

## Example Usage

```hcl
data "tencentcloud_sqlserver_account_db_attachments" "test" {
  instance_id  = tencentcloud_sqlserver_instance.test.id
  account_name = tencentcloud_sqlserver_account_db_attachment.test.account_name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) SQL Server instance ID that the account belongs to.
* `account_name` - (Optional) Name of the SQL Server account to be queried.
* `db_name` - (Optional) Name of the DB to be queried.
* `result_output_file` - (Optional) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of SQL Server account. Each element contains the following attributes:
  * `account_name` - SQL Server account name.
  * `db_name` - SQL Server DB name.
  * `instance_id` - SQL Server instance ID that the account belongs to.
  * `privilege` - Privilege of the account on DB. Valid value are `ReadOnly`, `ReadWrite`.


