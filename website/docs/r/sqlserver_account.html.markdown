---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_account"
sidebar_current: "docs-tencentcloud-resource-sqlserver_account"
description: |-
  Use this resource to create SQL Server account
---

# tencentcloud_sqlserver_account

Use this resource to create SQL Server account

## Example Usage

```hcl
resource "tencentcloud_sqlserver_account" "foo" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  name        = "tf_sqlserver_account"
  password    = "test1233"
  remark      = "testt"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) Instance ID that the account belongs to.
* `name` - (Required) Name of the SQL Server account.
* `password` - (Required) Password of the SQL Server account.
* `is_admin` - (Optional) Indicate that the account is root account or not.
* `remark` - (Optional) Remark of the SQL Server account.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Create time of the SQL Server account.
* `status` - Status of the SQL Server account. 1 for creating, 2 for running, 3 for modifying, 4 for resetting password, -1 for deleting.
* `update_time` - Last updated time of the SQL Server account.


## Import

SQL Server account can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_account.foo mssql-3cdq7kx5#tf_sqlserver_account
```

