---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_db"
sidebar_current: "docs-tencentcloud-resource-sqlserver_db"
description: |-
  Provides a SQL Server DB resource belongs to SQL Server instance.
---

# tencentcloud_sqlserver_db

Provides a SQL Server DB resource belongs to SQL Server instance.

## Example Usage

```hcl
resource "tencentcloud_sqlserver_db" "example" {
  instance_id = tencentcloud_sqlserver_instance.example.id
  name        = "example"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) SQLServer instance ID which DB belongs to.
* `name` - (Required, ForceNew) Name of SQL Server DB. The DataBase name must be unique and must be composed of numbers, letters and underlines, and the first one can not be underline.
* `charset` - (Optional, ForceNew) Character set DB uses. Valid values: `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `Chinese_PRC_BIN`, `Chinese_Taiwan_Stroke_CI_AS`, `SQL_Latin1_General_CP1_CI_AS`, and `SQL_Latin1_General_CP1_CS_AS`. Default value is `Chinese_PRC_CI_AS`.
* `remark` - (Optional) Remark of the DB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Database creation time.
* `status` - Database status, could be `creating`, `running`, `modifying` which means changing the remark, and `deleting`.


## Import

SQL Server DB can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_db.foo mssql-3cdq7kx5#db_name
```

