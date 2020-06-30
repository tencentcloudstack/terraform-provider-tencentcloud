---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_account_db_attachment"
sidebar_current: "docs-tencentcloud-resource-sqlserver_account_db_attachment"
description: |-
  Use this resource to create sqlserver account DB attachment
---

# tencentcloud_sqlserver_account_db_attachment

Use this resource to create sqlserver account DB attachment

## Example Usage

```hcl
resource "tencentcloud_sqlserver_account_db_attachment" "foo" {
  instance_id  = "mssql-3cdq7kx5"
  account_name = tencentcloud_sqlserver_account.example.name
  db_name      = tencentcloud_sqlserver_db.example.name
  privilege    = "ReadWrite"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required, ForceNew) SQL Server account name.
* `db_name` - (Required, ForceNew) SQL Server DB name.
* `instance_id` - (Required, ForceNew) SQL Server instance ID that the account belongs to.
* `privilege` - (Required) Privilege of the account on DB. Valid value are `READONLY`, `ReadWrite`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver account can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_account_db_attachment.foo mssql-3cdq7kx5#tf_sqlserver_account#test111
```

