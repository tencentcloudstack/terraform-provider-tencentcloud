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
  delete_subscribe_db    = false
  database_tuples {
    publish_database = tencentcloud_sqlserver_db.test_publish_subscribe.name
  }
}
```

## Argument Reference

The following arguments are supported:

* `database_tuples` - (Required) Database Publish and Publish relationship list. The elements inside can be deleted and added individually, but modification is not allowed.
* `publish_instance_id` - (Required, ForceNew) ID of the SQL Server instance which publish.
* `subscribe_instance_id` - (Required, ForceNew) ID of the SQL Server instance which subscribe.
* `delete_subscribe_db` - (Optional) Whether to delete the subscriber database when deleting the Publish and Subscribe. `true` for deletes the subscribe database, `false` for does not delete the subscribe database. default is `false`.
* `publish_subscribe_name` - (Optional) The name of the Publish and Subscribe. Default is `default_name`.

The `database_tuples` object supports the following:

* `publish_database` - (Required) Publish the database.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

SQL Server PublishSubscribe can be imported using the publish_sqlserver_id#subscribe_sqlserver_id, e.g.

```
$ terraform import tencentcloud_sqlserver_publish_subscribe.foo publish_sqlserver_id#subscribe_sqlserver_id
```

