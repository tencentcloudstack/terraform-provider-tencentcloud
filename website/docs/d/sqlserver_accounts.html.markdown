---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_accounts"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_accounts"
description: |-
  Use this data source to query the list of SQL Server accounts.
---

# tencentcloud_sqlserver_accounts

Use this data source to query the list of SQL Server accounts.

## Example Usage

### Pull instance account list

```hcl
data "tencentcloud_sqlserver_accounts" "example" {
  instance_id = "mssql-3cdq7kx5"
}
```

### Pull instance account list Filter by name

```hcl
data "tencentcloud_sqlserver_accounts" "example" {
  instance_id = "mssql-3cdq7kx5"
  name        = "myaccount"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) SQL server instance ID that the account belongs to.
* `name` - (Optional, String) Name of the SQL server account to be queried.
* `result_output_file` - (Optional, String) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - A list of SQL Server account. Each element contains the following attributes:
  * `create_time` - Create time of the SQL Server account.
  * `instance_id` - SQL server instance ID that the account belongs to.
  * `name` - Name of the SQL server account.
  * `remark` - Remark of the SQL Server account.
  * `status` - Status of the SQL Server account. `1` for creating, `2` for running, `3` for modifying, 4 for resetting password, -1 for deleting.
  * `update_time` - Last updated time of the SQL Server account.


