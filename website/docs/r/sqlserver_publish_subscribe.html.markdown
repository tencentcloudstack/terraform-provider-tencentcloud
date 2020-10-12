---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_publish_subscribe"
sidebar_current: "docs-tencentcloud-resource-sqlserver_publish_subscribe"
description: |-
  Provides a SQL Server PublishSubscribe resource belongs to SQL Server instance.
---

# tencentcloud_sqlserver_publish_subscribe

Provides a SQL Server PublishSubscribe resource belongs to SQL Server instance.

## Example Usage

```hcl
resource "tencentcloud_sqlserver_publish_subscribe" "example" {
  publish_instance_id    = tencentcloud_sqlserver_instance.publish_instance.id
  subscribe_instance_id  = tencentcloud_sqlserver_instance.subscribe_instance.id
  publish_subscribe_name = "example"
  database_tuples {
    publish_database   = tencentcloud_sqlserver_db.test_publish_subscribe.name
    subscribe_database = tencentcloud_sqlserver_db.test_publish_subscribe.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `database_tuples` - (Required) Database Publish and Publish relationship list. Modify database is not allowed.
* `publish_instance_id` - (Required, ForceNew) Publish the instance ID in the SQLServer instance.
* `subscribe_instance_id` - (Required, ForceNew) Subscribe the instance ID in the SQLServer instance.
* `publish_subscribe_name` - (Optional) The name of the Publish and Subscribe in the SQLServer instance. default is `default_name`.

The `database_tuples` object supports the following:

* `publish_database` - (Required) Publish the database.
* `subscribe_database` - (Required) Subscribe to the database.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SQL Server PublishSubscribe can be imported using the id, e.g.

```
$ terraform import tencentcloud_sqlserver_publish_subscribe.foo mssql-3cdq7kx5#db_name
```

