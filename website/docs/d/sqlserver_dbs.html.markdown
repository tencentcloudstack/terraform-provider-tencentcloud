---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_dbs"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_dbs"
description: |-
  Use this data source to query DB resources for the specific SQL Server instance.
---

# tencentcloud_sqlserver_dbs

Use this data source to query DB resources for the specific SQL Server instance.

## Example Usage

```hcl
data "tencentcloud_sqlserver_dbs" "example" {
  instance_id = "mssql-3cdq7kx5"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) SQL Server instance ID which DB belongs to.
* `result_output_file` - (Optional, String) Used to store results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `db_list` - A list of dbs belong to the specific instance. Each element contains the following attributes:
  * `charset` - Character set DB uses, could be `Chinese_PRC_CI_AS`, `Chinese_PRC_CS_AS`, `Chinese_PRC_BIN`, `Chinese_Taiwan_Stroke_CI_AS`, `SQL_Latin1_General_CP1_CI_AS`, and `SQL_Latin1_General_CP1_CS_AS`.
  * `create_time` - Database creation time.
  * `instance_id` - SQL Server instance ID which DB belongs to.
  * `name` - Name of DB.
  * `remark` - Remark of the DB.
  * `status` - Database status. Valid values are `creating`, `running`, `modifying`, `dropping`.


