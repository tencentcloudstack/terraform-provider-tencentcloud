---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_database_tde"
sidebar_current: "docs-tencentcloud-resource-sqlserver_database_tde"
description: |-
  Provides a resource to create a sqlserver database_tde
---

# tencentcloud_sqlserver_database_tde

Provides a resource to create a sqlserver database_tde

## Example Usage

```hcl
resource "tencentcloud_sqlserver_database_tde" "database_tde" {
  instance_id = "mssql-qelbzgwf"
  db_names    = ["keep_tde_db", "keep_tde_db2"]
  encryption  = "enable"
}
```

## Argument Reference

The following arguments are supported:

* `db_names` - (Required, Set: [`String`]) Database name list.
* `encryption` - (Required, String) `enable` - enable encryption, `disable` - disable encryption.
* `instance_id` - (Required, String) Instance ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

sqlserver database_tde can be imported using the id, e.g.

```
terraform import tencentcloud_sqlserver_database_tde.database_tde database_tde_id
```

