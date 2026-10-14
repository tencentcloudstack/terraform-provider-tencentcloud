---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_integration_task_node"
sidebar_current: "docs-tencentcloud-resource-wedata_integration_task_node"
description: |-
  Provides a resource to create a wedata integration_task_node
---

# tencentcloud_wedata_integration_task_node

Provides a resource to create a wedata integration_task_node

## Example Usage

```hcl
resource "tencentcloud_wedata_integration_task_node" "example" {
  project_id       = "1612982498218618880"
  task_id          = "20231022181114990"
  name             = "tf_example1"
  node_type        = "INPUT"
  data_source_type = "MYSQL"
  task_type        = 202
  task_mode        = 2
  node_info {
    datasource_id = "5085"
    config {
      name  = "Type"
      value = "MYSQL"
    }
    config {
      name  = "splitPk"
      value = "id"
    }
    config {
      name  = "PrimaryKey"
      value = "id"
    }
    config {
      name  = "isNew"
      value = "true"
    }
    config {
      name  = "PrimaryKey_INPUT_SYMBOL"
      value = "input"
    }
    config {
      name  = "splitPk_INPUT_SYMBOL"
      value = "input"
    }
    config {
      name  = "Database"
      value = "demo_mysql"
    }
    config {
      name  = "TableNames"
      value = "users"
    }
    config {
      name  = "SiblingNodes"
      value = "[]"
    }
    schema {
      id    = "471331072"
      name  = "id"
      type  = "INT"
      alias = "id"
    }
    schema {
      id    = "422052352"
      name  = "username"
      type  = "VARCHAR(50)"
      alias = "username"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `data_source_type` - (Required, String) Data source type: MYSQL, POSTGRE, ORACLE, SQLSERVER, FTP, HIVE, HDFS, ICEBERG, KAFKA, HBASE, SPARK, TBASE, DB2, DM, GAUSSDB, GBASE, IMPALA, ES, S3_DATAINSIGHT, GREENPLUM, PHOENIX, SAP_HANA, SFTP, OCEANBASE, CLICKHOUSE, KUDU, VERTICA, REDIS, COS, DLC, DORIS, CKAFKA, DTS_KAFKA, S3, CDW, TDSQLC, TDSQL, MONGODB, SYBASE, REST_API, StarRocks, TCHOUSE_X.
* `name` - (Required, String) Node Name.
* `node_info` - (Required, List) Node information.
* `node_type` - (Required, String) Node type: INPUT, OUTPUT, JOIN, FILTER, TRANSFORM.
* `project_id` - (Required, String) Project ID.
* `task_id` - (Required, String) The task id to which the node belongs.
* `task_mode` - (Required, Int) Task display mode, 0: canvas mode, 1: form mode.
* `task_type` - (Required, Int) Task type, 201: real-time task, 202: offline task.

The `config` object of `node_info` supports the following:

* `name` - (Optional, String) Configuration name.
* `value` - (Optional, String) Configuration value.

The `ext_config` object of `node_info` supports the following:

* `name` - (Optional, String) Configuration name.
* `value` - (Optional, String) Configuration value.

The `ext_config` object of `node_mapping` supports the following:

* `name` - (Optional, String) Configuration name.
* `value` - (Optional, String) Configuration value.

The `node_info` object supports the following:

* `app_id` - (Optional, String) User App Id.
* `config` - (Optional, List) Node configuration information.
* `create_time` - (Optional, String) Create time.
* `creator_uin` - (Optional, String) Creator User ID.
* `datasource_id` - (Optional, String) Datasource ID.
* `ext_config` - (Optional, List) Node extension configuration information.
* `node_mapping` - (Optional, List) Node mapping.
* `operator_uin` - (Optional, String) Operator User ID.
* `owner_uin` - (Optional, String) Owner User ID.
* `schema` - (Optional, List) Schema information.
* `update_time` - (Optional, String) Update time.

The `node_mapping` object of `node_info` supports the following:

* `ext_config` - (Optional, List) Node extension configuration information.
* `schema_mappings` - (Optional, List) Schema mapping information.
* `sink_id` - (Optional, String) Sink node ID.
* `source_id` - (Optional, String) Source node ID.
* `source_schema` - (Optional, List) Source node schema information.

The `properties` object of `schema` supports the following:

* `name` - (Optional, String) Attributes name.
* `value` - (Optional, String) Attributes value.

The `properties` object of `source_schema` supports the following:

* `name` - (Optional, String) Attributes name.
* `value` - (Optional, String) Attributes value.

The `schema_mappings` object of `node_mapping` supports the following:

* `sink_schema_id` - (Required, String) Schema ID from sink node.
* `source_schema_id` - (Required, String) Schema ID from source node.

The `schema` object of `node_info` supports the following:

* `id` - (Required, String) Schema ID.
* `name` - (Required, String) Schema name.
* `type` - (Required, String) Schema type.
* `alias` - (Optional, String) Schema alias.
* `comment` - (Optional, String) Schema comment.
* `properties` - (Optional, List) Schema extended attributes.
* `value` - (Optional, String) Schema value.

The `source_schema` object of `node_mapping` supports the following:

* `id` - (Required, String) Schema ID.
* `name` - (Required, String) Schema name.
* `type` - (Required, String) Schema type.
* `alias` - (Optional, String) Schema alias.
* `comment` - (Optional, String) Schema comment.
* `properties` - (Optional, List) Schema extended attributes.
* `value` - (Optional, String) Schema value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `node_id` - Node ID.


