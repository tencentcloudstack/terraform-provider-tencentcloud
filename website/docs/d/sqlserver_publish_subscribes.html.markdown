---
subcategory: "SQLServer"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_sqlserver_publish_subscribes"
sidebar_current: "docs-tencentcloud-datasource-sqlserver_publish_subscribes"
description: |-
  Use this data source to query Publish Subscribe resources for the specific SQL Server instance.
---

# tencentcloud_sqlserver_publish_subscribes

Use this data source to query Publish Subscribe resources for the specific SQL Server instance.

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

* `instance_id` - (Required) ID of the SQL Server instance.
* `pub_or_sub_instance_id` - (Optional) The subscribe/publish instance ID. It is related to whether the `instance_id` is a publish instance or a subscribe instance. when `instance_id` is a publish instance, this field is filtered according to the subscribe instance ID; when `instance_id` is a subscribe instance, this field is filtering according to the publish instance ID.
* `pub_or_sub_instance_ip` - (Optional) The intranet IP of the subscribe/publish instance. It is related to whether the `instance_id` is a publish instance or a subscribe instance. when `instance_id` is a publish instance, this field is filtered according to the intranet IP of the subscribe instance; when `instance_id` is a subscribe instance, this field is based on the publish instance intranet IP filter.
* `publish_database` - (Optional) Name of publish database.
* `publish_subscribe_id` - (Optional) The id of the Publish and Subscribe.
* `publish_subscribe_name` - (Optional) The name of the Publish and Subscribe.
* `result_output_file` - (Optional) Used to store results.
* `subscribe_database` - (Optional) Name of subscribe database.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `publish_subscribe_list` - Publish and subscribe list. Each element contains the following attributes.
  * `database_tuples` - Database Publish and Publish relationship list.
    * `last_sync_time` - Last sync time.
    * `publish_database` - Name of the publish SQL Server instance.
    * `status` - Publish and subscribe status between databases, valid values are `running`, `success`, `fail`, `unknow`.
    * `subscribe_database` - Name of the subscribe SQL Server instance.
  * `publish_instance_id` - ID of the SQL Server instance which publish.
  * `publish_instance_ip` - IP of the the SQL Server instance which publish.
  * `publish_instance_name` - Name of the SQL Server instance which publish.
  * `publish_subscribe_id` - The id of the Publish and Subscribe.
  * `publish_subscribe_name` - The name of the Publish and Subscribe.
  * `subscribe_instance_id` - ID of the SQL Server instance which subscribe.
  * `subscribe_instance_ip` - IP of the SQL Server instance which subscribe.
  * `subscribe_instance_name` - Name of the SQL Server instance which subscribe.


