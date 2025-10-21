---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_connect_resource"
sidebar_current: "docs-tencentcloud-resource-ckafka_connect_resource"
description: |-
  Provides a resource to create a ckafka connect_resource
---

# tencentcloud_ckafka_connect_resource

Provides a resource to create a ckafka connect_resource

## Example Usage

```hcl
resource "tencentcloud_ckafka_connect_resource" "connect_resource" {
  resource_name = "terraform-test"
  type          = "MYSQL"
  description   = "for terraform test"

  mysql_connect_param {
    port        = 3306
    user_name   = "root"
    password    = "xxxxxxxxx"
    resource    = "cdb-fitq5t9h"
    service_vip = "172.16.80.59"
    uniq_vpc_id = "vpc-4owdpnwr"
    self_built  = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `resource_name` - (Required, String) connection source name.
* `type` - (Required, String, ForceNew) connection source type.
* `clickhouse_connect_param` - (Optional, List) ClickHouse configuration, required when Type is CLICKHOUSE.
* `description` - (Optional, String) Connection source description.
* `doris_connect_param` - (Optional, List) Doris configuration, required when Type is DORIS.
* `dts_connect_param` - (Optional, List) Dts configuration, required when Type is DTS.
* `es_connect_param` - (Optional, List) Es configuration, required when Type is ES.
* `kafka_connect_param` - (Optional, List) Kafka configuration, required when Type is KAFKA.
* `mariadb_connect_param` - (Optional, List) Maria DB configuration, required when Type is MARIADB.
* `mongodb_connect_param` - (Optional, List) Mongo DB configuration, required when Type is MONGODB.
* `mysql_connect_param` - (Optional, List) MySQL configuration, required when Type is MYSQL or TDSQL C_MYSQL.
* `postgresql_connect_param` - (Optional, List) Postgresql configuration, required when Type is POSTGRESQL or TDSQL C POSTGRESQL.
* `sqlserver_connect_param` - (Optional, List) SQLServer configuration, required when Type is SQLSERVER.

The `clickhouse_connect_param` object supports the following:

* `password` - (Required, String) Password for Clickhouse connection source.
* `port` - (Required, Int) Clickhouse connection port.
* `resource` - (Required, String) Instance resources for Click House connection sources.
* `self_built` - (Required, Bool) Whether the Clickhouse connection source is a self-built cluster.
* `user_name` - (Required, String) The username of the clickhouse connection source.
* `is_update` - (Optional, Bool) Whether to update to the associated Datahub task, default: false.
* `service_vip` - (Optional, String) Instance VIP of the ClickHouse connection source, when it is a Tencent Cloud instance, it is required.
* `uniq_vpc_id` - (Optional, String) The vpc Id of the source of the ClickHouse connection, when it is a Tencent Cloud instance, it is required.

The `doris_connect_param` object supports the following:

* `password` - (Required, String) Doris  password.
* `port` - (Required, Int) Doris jdbc CLB port, Usually mapped to port 9030 of fe.
* `resource` - (Required, String) Doris  instanceId.
* `user_name` - (Required, String) Doris  The username of the connection source.
* `be_port` - (Optional, Int) Doris http CLB port, Usually mapped to port 8040 of be.
* `is_update` - (Optional, Bool) Whether to update to the associated Datahub task, default: false.
* `self_built` - (Optional, Bool) Doris Whether the connection source is a self-built cluster, default: false.
* `service_vip` - (Optional, String) Doris vip, When it is a Tencent Cloud instance, it is required.
* `uniq_vpc_id` - (Optional, String) Doris vpcId, When it is a Tencent Cloud instance, it is required.

The `dts_connect_param` object supports the following:

* `group_id` - (Required, String) Id of the Dts consumption group.
* `password` - (Required, String) The password of the Dts consumption group.
* `port` - (Required, Int) Dts port.
* `resource` - (Required, String) Dts instance Id.
* `topic` - (Required, String) Topic subscribed by Dts.
* `user_name` - (Required, String) The account number of the Dts consumption group.
* `is_update` - (Optional, Bool) Whether to update to the associated Datahub task, default: false.

The `es_connect_param` object supports the following:

* `password` - (Required, String) Es The password of the connection source.
* `port` - (Required, Int) Es port.
* `resource` - (Required, String) Instance resource of Es connection source.
* `self_built` - (Required, Bool) Whether the Es connection source is a self-built cluster.
* `user_name` - (Required, String) Es The username of the connection source.
* `is_update` - (Optional, Bool) Whether to update to the associated Datahub task, default: false.
* `service_vip` - (Optional, String) The instance vip of the Es connection source, when it is a Tencent Cloud instance, it is required.
* `uniq_vpc_id` - (Optional, String) The vpc Id of the Es connection source, when it is a Tencent Cloud instance, it is required.

The `kafka_connect_param` object supports the following:

* `broker_address` - (Optional, String) Kafka broker ip, Mandatory when self-built.
* `is_update` - (Optional, Bool) Whether to update to the associated Dip task, default: false.
* `region` - (Optional, String) CKafka instanceId region, Required when crossing regions.
* `resource` - (Optional, String) Kafka instanceId, When it is a Tencent Cloud instance, it is required.
* `self_built` - (Optional, Bool) Whether it is a self-built cluster, default: false.

The `mariadb_connect_param` object supports the following:

* `password` - (Required, String) MariaDB password.
* `port` - (Required, Int) MariaDB port.
* `resource` - (Required, String) MariaDB instanceId.
* `user_name` - (Required, String) MariaDB The username of the connection source.
* `is_update` - (Optional, Bool) Whether to update to the associated Datahub task, default: false.
* `service_vip` - (Optional, String) The instance vip of the Maria DB connection source, when it is a Tencent Cloud instance, it is required.
* `uniq_vpc_id` - (Optional, String) MariaDB vpcId, When it is a Tencent Cloud instance, it is required.

The `mongodb_connect_param` object supports the following:

* `password` - (Required, String) Password for the source of the Mongo DB connection.
* `port` - (Required, Int) MongoDB port.
* `resource` - (Required, String) Instance resource of Mongo DB connection source.
* `self_built` - (Required, Bool) Whether the Mongo DB connection source is a self-built cluster.
* `user_name` - (Required, String) The username of the Mongo DB connection source.
* `is_update` - (Optional, Bool) Whether to update to the associated Datahub task, default: false.
* `service_vip` - (Optional, String) The instance VIP of the Mongo DB connection source, when it is a Tencent Cloud instance, it is required.
* `uniq_vpc_id` - (Optional, String) The vpc Id of the Mongo DB connection source, which is required when it is a Tencent Cloud instance.

The `mysql_connect_param` object supports the following:

* `password` - (Required, String) Mysql connection source password.
* `port` - (Required, Int) MySQL port.
* `resource` - (Required, String) Instance resource of My SQL connection source.
* `user_name` - (Required, String) Username of Mysql connection source.
* `cluster_id` - (Optional, String) Required when type is TDSQL C_MYSQL.
* `is_update` - (Optional, Bool) Whether to update to the associated Datahub task, default: false.
* `self_built` - (Optional, Bool) Mysql Whether the connection source is a self-built cluster, default: false.
* `service_vip` - (Optional, String) The instance vip of the MySQL connection source, when it is a Tencent Cloud instance, it is required.
* `uniq_vpc_id` - (Optional, String) The vpc Id of the My SQL connection source, when it is a Tencent Cloud instance, it is required.

The `postgresql_connect_param` object supports the following:

* `password` - (Required, String) PostgreSQL password.
* `port` - (Required, Int) PostgreSQL port.
* `resource` - (Required, String) PostgreSQL instanceId.
* `user_name` - (Required, String) PostgreSQL The username of the connection source.
* `cluster_id` - (Optional, String) Required when type is TDSQL C_POSTGRESQL.
* `is_update` - (Optional, Bool) Whether to update to the associated Datahub task, default: false.
* `self_built` - (Optional, Bool) PostgreSQL Whether the connection source is a self-built cluster, default: false.
* `service_vip` - (Optional, String) The instance VIP of the Postgresql connection source, when it is a Tencent Cloud instance, it is required.
* `uniq_vpc_id` - (Optional, String) The instance vpcId of the Postgresql connection source, when it is a Tencent Cloud instance, it is required.

The `sqlserver_connect_param` object supports the following:

* `password` - (Required, String) SQLServer password.
* `port` - (Required, Int) SQLServer port.
* `resource` - (Required, String) SQLServer instanceId.
* `user_name` - (Required, String) SQLServer The username of the connection source.
* `is_update` - (Optional, Bool) Whether to update to the associated Dip task, default: false.
* `service_vip` - (Optional, String) SQLServer instance vip, When it is a Tencent Cloud instance, it is required.
* `uniq_vpc_id` - (Optional, String) SQLServer vpcId, When it is a Tencent Cloud instance, it is required.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ckafka connect_resource can be imported using the id, e.g.

```
terraform import tencentcloud_ckafka_connect_resource.connect_resource connect_resource_id
```

