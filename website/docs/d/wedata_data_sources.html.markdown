---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_data_sources"
sidebar_current: "docs-tencentcloud-datasource-wedata_data_sources"
description: |-
  Use this data source to query detailed information of WeData data sources
---

# tencentcloud_wedata_data_sources

Use this data source to query detailed information of WeData data sources

## Example Usage

```hcl
data "tencentcloud_wedata_data_sources" "example" {
  project_id   = "2982667120655491072"
  name         = "tf_example"
  display_name = "display_name"
  type         = ["MYSQL", "ORACLE"]
  creator      = "user"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, String) Project ID.
* `creator` - (Optional, String) Creator.
* `display_name` - (Optional, String) Data source display name.
* `name` - (Optional, String) Data source name.
* `result_output_file` - (Optional, String) Used to save results.
* `type` - (Optional, Set: [`String`]) Data source type: enumeration values.

- MYSQL
- TENCENT_MYSQL
- POSTGRE
- ORACLE
- SQLSERVER
- FTP
- HIVE
- HUDI
- HDFS
- ICEBERG
- KAFKA
- HBASE
- SPARK
- VIRTUAL
- TBASE
- DB2
- DM
- GAUSSDB
- GBASE
- IMPALA
- ES
- TENCENT_ES
- GREENPLUM
- PHOENIX
- SAP_HANA
- SFTP
- OCEANBASE
- CLICKHOUSE
- KUDU
- VERTICA
- REDIS
- COS
- DLC
- DORIS
- CKAFKA
- S3
- TDSQL
- TDSQL_MYSQL
- MONGODB
- TENCENT_MONGODB
- REST_API
- SuperSQL
- PRESTO
- TiDB
- StarRocks
- Trino
- Kyuubi
- TCHOUSE_X
- TCHOUSE_P
- TCHOUSE_C
- TCHOUSE_D
- INFLUXDB
- BIG_QUERY
- SSH
- BLOB.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Data source list.
  * `category` - Data source category:

- DB - custom source
- CLUSTER - system source.
  * `create_time` - Time.
  * `create_user` - Data source creator.
  * `description` - Data source description information.
  * `dev_con_properties` - Same as params, contains data for development data source.
  * `display_name` - Data source display name, for visual viewing.
  * `id` - Data source ID.
  * `modify_time` - Modification time.
  * `modify_user` - Modifier.
  * `name` - Data source name.
  * `prod_con_properties` - Data source configuration information, stored in JSON KV format, varies by data source type.
  * `project_id` - Belonging project ID.
  * `project_name` - Belonging project name.
  * `type` - Data source type: enumeration values.


