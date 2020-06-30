---
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
data "tencentcloud_sqlserver_accounts" "account" {
  instance_id  = "mssql-3cdq7kx5"
  account_name = "myaccount"
}

data "tencentcloud_sqlserver_accounts" "db" {
  instance_id = "mssql-3cdq7kx5"
  db_name     = "mydb"
}

data "tencentcloud_sqlserver_accounts" "instance" {
  instance_id = "mssql-3cdq7kx5"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) SQL Server instance ID that the account belongs to.
* `account_name` - (Optional) Name of the SQL Server account to be queried.
* `db_name` - (Optional) Name of the SQL Server account to be queried.
* `result_output_file` - (Optional) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of SQL Server account. Each element contains the following attributes:
  * `account_name` - SQL Server account name.
  * `db_name` - SQL Server DB name.
  * `instance_id` - SQL Server instance ID that the account belongs to.
  * `privilege` - Privilege of the account on DB. Valid value are `READONLY`, `ReadWrite`.


