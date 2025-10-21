---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_connect_resource"
sidebar_current: "docs-tencentcloud-datasource-ckafka_connect_resource"
description: |-
  Use this data source to query detailed information of ckafka connect_resource
---

# tencentcloud_ckafka_connect_resource

Use this data source to query detailed information of ckafka connect_resource

## Example Usage

```hcl
data "tencentcloud_ckafka_connect_resource" "connect_resource" {
}
```

## Argument Reference

The following arguments are supported:

* `limit` - (Optional, Int) Return the number, the default is 20, the maximum is 100.
* `offset` - (Optional, Int) Page offset, default is 0.
* `resource_region` - (Optional, String) Keyword query of the connection source, query the connection in the connection management list in the local region according to the region (only support the connection source containing the region input).
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) Keyword for search.
* `type` - (Optional, String) connection source type.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `result` - Connection source list.
  * `connect_resource_list` - Resource List.
    * `clickhouse_connect_param` - ClickHouse configuration, returned when Type is CLICKHOUSE.
      * `is_update` - Whether to update to the associated Datahub task.
      * `password` - The password of the connection source.
      * `port` - ClickHouse port.
      * `resource` - Instance resource of connection source.
      * `self_built` - Whether the connection source is a self-built cluster.
      * `service_vip` - Instance VIP of the connection source, when it is a Tencent Cloud instance, it is required.
      * `uniq_vpc_id` - The vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.
      * `user_name` - The username of the connection source.
    * `create_time` - Creation time.
    * `ctsdb_connect_param` - Ctsdb configuration, returned when Type is CTSDB.
      * `password` - The password of the connection source.
      * `port` - Ctsdb port.
      * `resource` - Instance resource of connection source.
      * `service_vip` - Ctsdb vip.
      * `uniq_vpc_id` - Ctsdb vpcId.
      * `user_name` - The username of the connection source.
    * `current_step` - The current step of the connection source.
    * `datahub_task_count` - The number of Datahub tasks associated with this connection source.
    * `description` - Description.
    * `doris_connect_param` - Doris Configuration, returned when Type is DORIS.
      * `be_port` - Doris's http load balancing connection port, usually mapped to be's 8040 port.
      * `is_update` - Whether to update to the associated Datahub task.
      * `password` - The password of the connection source.
      * `port` - Doris jdbc Load balancing connection port, usually mapped to port 9030 of fe.
      * `resource` - Instance resource of connection source.
      * `self_built` - Whether the connection source is a self-built cluster.
      * `service_vip` - Instance VIP of the connection source, when it is a Tencent Cloud instance, it is required.
      * `uniq_vpc_id` - The vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.
      * `user_name` - The username of the connection source.
    * `dts_connect_param` - Dts configuration, returned when Type is DTS.
      * `group_id` - The id of the Dts consumer group.
      * `is_update` - Whether to update to the associated Datahub task.
      * `password` - The password of the Dts consumer group.
      * `port` - Dts port.
      * `resource` - Dts Id.
      * `topic` - Topic subscribed by Dts.
      * `user_name` - The UserName of the Dts consumer group.
    * `error_message` - Error Messages.
    * `es_connect_param` - Es configuration, return when Type is ES.
      * `is_update` - Whether to update to the associated Datahub task.
      * `password` - The password of the connection source.
      * `port` - ES port.
      * `resource` - Instance resource of connection source.
      * `self_built` - Whether the connection source is a self-built cluster.
      * `service_vip` - Instance VIP of the connection source, when it is a Tencent Cloud instance, it is required.
      * `uniq_vpc_id` - The vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.
      * `user_name` - The username of the connection source.
    * `kafka_connect_param` - Kafka configuration, returned when Type is KAFKA.
      * `broker_address` - Broker address for Kafka connection, required for self-build.
      * `is_update` - Whether to update to the associated Dip task.
      * `region` - Instance resource region of CKafka connection source, required when crossing regions.
      * `resource` - Instance resource of Kafka connection source, required when not self-built.
      * `self_built` - Whether it is a self-built cluster.
    * `maria_db_connect_param` - Mariadb configuration, returned when Type is MARIADB.
      * `is_update` - Whether to update to the associated Datahub task.
      * `password` - The password of the connection source.
      * `port` - MariaDB port.
      * `resource` - Instance resource of connection source.
      * `service_vip` - Instance VIP of the connection source, when it is a Tencent Cloud instance, it is required.
      * `uniq_vpc_id` - The vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.
      * `user_name` - The username of the connection source.
    * `mongo_db_connect_param` - Mongo DB configuration, returned when Type is MONGODB.
      * `is_update` - Whether to update to the associated Datahub task.
      * `password` - The password of the connection source.
      * `port` - MongoDB port.
      * `resource` - Instance resource of connection source.
      * `self_built` - Whether the connection source is a self-built cluster.
      * `service_vip` - Instance VIP of the connection source, when it is a Tencent Cloud instance, it is required.
      * `uniq_vpc_id` - The vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.
      * `user_name` - The username of the connection source.
    * `mysql_connect_param` - Mysql configuration, returned when Type is MYSQL or TDSQL C MYSQL.
      * `cluster_id` - Required when type is TDSQL C_MYSQL.
      * `is_update` - Whether to update to the associated Datahub task.
      * `password` - The password of the connection source.
      * `port` - MySQL port.
      * `resource` - MySQL Instance resource of connection source.
      * `self_built` - Mysql Whether the connection source is a self-built cluster.
      * `service_vip` - Instance VIP of the connection source, when it is a Tencent Cloud instance, it is required.
      * `uniq_vpc_id` - The vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.
      * `user_name` - The username of the connection source.
    * `postgre_sql_connect_param` - Postgresql configuration, returned when Type is POSTGRESQL or TDSQL C POSTGRESQL.
      * `cluster_id` - Required when type is TDSQL C_POSTGRESQL.
      * `is_update` - Whether to update to the associated Datahub task.
      * `password` - The password of the connection source.
      * `port` - PostgreSQL port.
      * `resource` - Instance resource of connection source.
      * `self_built` - Whether the connection source is a self-built cluster.
      * `service_vip` - Instance VIP of the connection source, when it is a Tencent Cloud instance, it is required.
      * `uniq_vpc_id` - The vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.
      * `user_name` - The username of the connection source.
    * `resource_id` - Resource id.
    * `resource_name` - Resource name.
    * `sql_server_connect_param` - SQL Server configuration, returned when Type is SQLSERVER.
      * `is_update` - Whether to update to the associated Dip task.
      * `password` - The password of the connection source.
      * `port` - SQLServer port.
      * `resource` - Instance resource of connection source.
      * `service_vip` - Instance VIP of the connection source, when it is a Tencent Cloud instance, it is required.
      * `uniq_vpc_id` - The vpc Id of the connection source, when it is a Tencent Cloud instance, it is required.
      * `user_name` - The username of the connection source.
    * `status` - Resource status.
    * `step_list` - Step List.
    * `task_progress` - Creation progress percentage.
    * `type` - Resource type.
  * `total_count` - Number of connection sources.


