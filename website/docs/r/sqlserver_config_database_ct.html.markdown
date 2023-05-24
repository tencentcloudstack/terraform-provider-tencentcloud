---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_config_database_ct"
sidebar_current: "docs-tencentcloud-resource-sqlserver_config_database_ct"
description: |-
  Provides a resource to create a sqlserver config_database_ct
---

# tencentcloud_sqlserver_config_database_ct

Provides a resource to create a sqlserver config_database_ct

## Example Usage

```hcl
resource "tencentcloud_sqlserver_config_database_ct" "config_database_ct" {
  db_name              = "keep_pubsub_db2"
  modify_type          = "disable"
  instance_id          = "mssql-qelbzgwf"
  change_retention_day = 7
}

resource "tencentcloud_sqlserver_config_database_ct" "config_database_ct" {
  db_name     = "keep_pubsub_db2"
  modify_type = "disable"
  instance_id = "mssql-qelbzgwf"
}
```

## Argument Reference

The following arguments are supported:

* `db_name` - (Required, String) database name.
* `instance_id` - (Required, String) Instance ID.
* `modify_type` - (Required, String) Enable or disable CT. Valid values: enable, disable.
* `change_retention_day` - (Optional, Int) Retention period (in days) of change tracking information when CT is enabled. Value range: 3-30. Default value: 3.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver tencentcloud_sqlserver_config_database_ct can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_database_ct.config_database_ct config_database_ct_id
```

