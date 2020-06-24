---
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_db"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_db"
description: |-
  Use this data source to query DB resources for the specific SQLServer instance.
---

# tencentcloud_sqlserver_db

Use this data source to query DB resources for the specific SQLServer instance.

## Example Usage

```hcl
resource "tencentcloud_sqlserver_db" "mysqlserver_db" {
  instance_id = "mssql-XXXXXX"
  name        = "sqlserver_db_terraform"
  charset     = "Chinese_PRC_BIN"
  remark      = "test-remark"
}

data "tencentcloud_sqlserver_db" "mysqlserver" {
  instance_id = tencentcloud_sqlserver_db.mysqlserver_db.instance_id
  name        = tencentcloud_sqlserver_db.mysqlserver_db.name
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) SQLServer instance ID which DB belongs to.
* `name` - (Required) Name of DB.
* `result_output_file` - (Optional) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `charset` - Character set DB uses, could be `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `Chinese_PRC_BIN`, `Chinese_Taiwan_Stroke_CI_AS`, `SQL_Latin1_General_CP1_CI_AS`, and `SQL_Latin1_General_CP1_CS_AS`. Default value is `Chinese_PRC_CI_AS`.
* `create_time` - Database creation time.
* `remark` - Remark of the DB.
* `status` - Database status. Valid values are `creating`, `running`, `modifying`, `dropping`.


