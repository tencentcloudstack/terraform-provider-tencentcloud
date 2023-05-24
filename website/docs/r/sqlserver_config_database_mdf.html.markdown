---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_config_database_mdf"
sidebar_current: "docs-tencentcloud-resource-sqlserver_config_database_mdf"
description: |-
  Provides a resource to create a sqlserver config_database_mdf
---

# tencentcloud_sqlserver_config_database_mdf

Provides a resource to create a sqlserver config_database_mdf

## Example Usage

```hcl
resource "tencentcloud_sqlserver_config_database_mdf" "config_database_mdf" {
  db_name     = "keep_pubsub_db2"
  instance_id = "mssql-qelbzgwf"
}
```

## Argument Reference

The following arguments are supported:

* `db_name` - (Required, String) Array of database names.
* `instance_id` - (Required, String) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver config_database_mdf can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_config_database_mdf.config_database_mdf config_database_mdf_id
```

