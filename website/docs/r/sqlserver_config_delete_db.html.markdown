---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_config_delete_db"
sidebar_current: "docs-tencentcloud-resource-sqlserver_config_delete_db"
description: |-
  Provides a resource to create a sqlserver config_delete_db
---

# tencentcloud_sqlserver_config_delete_db

Provides a resource to create a sqlserver config_delete_db

## Example Usage

```hcl
resource "tencentcloud_sqlserver_config_delete_db" "config_delete_db" {
  instance_id = "mssql-i1z41iwd"
  name        =
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID.
* `name` - (Required, String) collection of database name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



