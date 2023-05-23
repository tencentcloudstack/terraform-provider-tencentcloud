---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_config_database_cdc"
sidebar_current: "docs-tencentcloud-resource-sqlserver_config_database_cdc"
description: |-
  Provides a resource to create a sqlserver config_database_cdc
---

# tencentcloud_sqlserver_config_database_cdc

Provides a resource to create a sqlserver config_database_cdc

## Example Usage

```hcl
resource "tencentcloud_sqlserver_config_database_cdc" "config_database_cdc" {
  db_name     = "keep_pubsub_db2"
  modify_type = "disable"
  instance_id = "mssql-qelbzgwf"
}
```

## Argument Reference

The following arguments are supported:

* `db_name` - (Required, String) database name.
* `instance_id` - (Required, String) Instance ID.
* `modify_type` - (Required, String) Enable or disable CDC. Valid values: enable, disable.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver config_database_cdc can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_database_cdc.config_database_cdc config_database_cdc_id
```

