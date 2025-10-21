---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_datahub_task"
sidebar_current: "docs-tencentcloud-datasource-ckafka_datahub_task"
description: |-
  Use this data source to query detailed information of ckafka datahub_task
---

# tencentcloud_ckafka_datahub_task

Use this data source to query detailed information of ckafka datahub_task

## Example Usage

```hcl
data "tencentcloud_ckafka_datahub_task" "datahub_task" {
}
```

## Argument Reference

The following arguments are supported:

* `resource` - (Optional, String) Resource.
* `result_output_file` - (Optional, String) Used to save results.
* `search_word` - (Optional, String) search key.
* `source_type` - (Optional, String) The source type.
* `target_type` - (Optional, String) Destination type of dump.
* `task_type` - (Optional, String) Task type, SOURCE|SINK.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `task_list` - Datahub task information list.
  * `create_time` - CreateTime.
  * `datahub_id` - Datahub Id.
  * `error_message` - ErrorMessage.
  * `source_resource` - data resource.
    * `click_house_param` - ClickHouse config, Type CLICKHOUSE requierd.
      * `cluster` - ClickHouse cluster.
      * `database` - ClickHouse database name.
      * `drop_cls` - When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.
        * `drop_cls_log_set` - cls LogSet id.
        * `drop_cls_owneruin` - cls account.
        * `drop_cls_region` - cls region.
        * `drop_cls_topic_id` - cls topicId.
        * `drop_invalid_message_to_cls` - Whether to deliver to cls.
      * `drop_invalid_message` - Whether ClickHouse discards the message that fails to parse, the default is true.
      * `ip` - ClickHouse ip.
      * `password` - ClickHouse passwd.
      * `port` - ClickHouse port.
      * `resource` - resource id.
      * `schema` - ClickHouse schema.
        * `allow_null` - Whether the column item is allowed to be empty.
        * `column_name` - column name.
        * `json_key` - The json Key name corresponding to this column.
        * `type` - type of table column.
      * `self_built` - Whether it is a self-built cluster.
      * `service_vip` - instance vip.
      * `table` - ClickHouse table.
      * `type` - ClickHouse type, emr-clickhouse: emr; cdw-clickhouse: cdwch; selfBuilt: "".
      * `uniq_vpc_id` - instance vpc id.
      * `user_name` - ClickHouse user name.
    * `cls_param` - Cls configuration, Required when Type is CLS.
      * `content_key` - Required when Decode Json is false.
      * `decode_json` - Whether the produced information is in json format.
      * `log_set` - LogSet id.
      * `resource` - cls id.
      * `time_field` - Specify the content of a field in the message as the time of the cls log. The format of the field content needs to be a second-level timestamp.
    * `cos_param` - Cos configuration, required when Type is COS.
      * `aggregate_batch_size` - The size of aggregated messages MB.
      * `aggregate_interval` - time interval.
      * `bucket_name` - cos bucket name.
      * `directory_time_format` - Partition format formatted according to strptime time.
      * `format_output_type` - The file format after message aggregation csv|json.
      * `object_key_prefix` - Dumped object directory prefix.
      * `object_key` - ObjectKey.
      * `region` - region code.
    * `ctsdb_param` - Ctsdb configuration, Required when Type is CTSDB.
      * `ctsdb_metric` - Ctsdb metric.
      * `resource` - resource id.
    * `dts_param` - Dts configuration, required when Type is DTS.
      * `group_id` - Dts consumer group Id.
      * `group_password` - Dts consumer group passwd.
      * `group_user` - Dts account.
      * `ip` - Dts connection ip.
      * `port` - Dts connection port.
      * `resource` - Dts instance Id.
      * `topic` - Dts topic.
      * `tran_sql` - False to synchronize the original data, true to synchronize the parsed json format data, the default is true.
    * `es_param` - Es configuration, required when Type is ES.
      * `content_key` - key for data in non-json format.
      * `database_primary_key` - When the message dumped to ES is the binlog of Database, if you need to synchronize database operations, that is, fill in the primary key of the database table when adding, deleting, and modifying operations to ES.
      * `date_format` - Es date suffix.
      * `document_id_field` - The field name of the document ID value dumped into Es.
      * `drop_cls` - When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.
        * `drop_cls_log_set` - cls log set.
        * `drop_cls_owneruin` - Delivery account of cls.
        * `drop_cls_region` - The region where the cls is delivered.
        * `drop_cls_topic_id` - topic of cls.
        * `drop_invalid_message_to_cls` - Whether to deliver to cls.
      * `drop_dlq` - dead letter queue.
        * `dlq_type` - dlq type, CKAFKA|TOPIC.
        * `kafka_param` - Ckafka type dlq.
          * `compression_type` - Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.
          * `connector_sync_type` - ConnectorSyncType.
          * `enable_toleration` - Enable the fault-tolerant instance and enable the dead-letter queue.
          * `keep_partition` - KeepPartition.
          * `msg_multiple` - 1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).
          * `offset_type` - Offset type, initial position earliest, latest position latest, time point position timestamp.
          * `partition_num` - Partition num.
          * `qps_limit` - Qps limit.
          * `resource_name` - resource id name.
          * `resource` - resource id.
          * `self_built` - Whether it is a self-built cluster.
          * `start_time` - It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.
          * `table_mappings` - The route from Table to Topic must be passed when the Distribute to multiple topics switch is turned on.
            * `database` - database name.
            * `table` - Table name, multiple tables, separated by (commas).
            * `topic_id` - Topic ID.
            * `topic` - Topic name.
          * `topic_id` - Topic Id.
          * `topic` - Topic name, multiple separated by,.
          * `use_auto_create_topic` - whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).
          * `use_table_mapping` - Distribute to multiple topics switch, the default is false.
          * `zone_id` - Zone ID.
        * `max_retry_attempts` - retry times.
        * `retry_interval` - retry interval.
        * `topic_param` - DIP Topic type dead letter queue.
          * `compression_type` - Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.
          * `msg_multiple` - 1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).
          * `offset_type` - Offset type, initial position earliest, latest position latest, time point position timestamp.
          * `resource` - The topic name of the topic sold separately.
          * `start_time` - It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.
          * `topic_id` - TopicId.
          * `use_auto_create_topic` - whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).
        * `type` - type, DLQ dead letter queue, IGNORE_ERROR|DROP.
      * `drop_invalid_json_message` - Whether Es discards messages in non-json format.
      * `drop_invalid_message` - Whether Es discards the message of parsing failure.
      * `index_type` - Es custom index name type, STRING, JSONPATH, the default is STRING.
      * `index` - Es index name.
      * `password` - Es Password.
      * `port` - Es connection port.
      * `resource` - Resource.
      * `self_built` - Whether it is a self-built cluster.
      * `service_vip` - instance vip.
      * `uniq_vpc_id` - instance vpc id.
      * `user_name` - Es UserName.
    * `event_bus_param` - EB configuration, required when type is EB.
      * `function_name` - SCF function name.
      * `namespace` - SCF namespace.
      * `qualifier` - SCF version and alias.
      * `resource` - instance id.
      * `self_built` - Whether it is a self-built cluster.
      * `type` - resource type EB_COS/EB_ES/EB_CLS.
    * `kafka_param` - ckafka configuration, required when Type is KAFKA.
      * `compression_type` - Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.
      * `connector_sync_type` - ConnectorSyncType.
      * `enable_toleration` - enable dead letter queue.
      * `keep_partition` - KeepPartition.
      * `msg_multiple` - 1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).
      * `offset_type` - Offset type, from beginning:earliest, from latest:latest, from specific time:timestamp.
      * `partition_num` - the partition num of the topic.
      * `qps_limit` - Qps(query per seconds) limit.
      * `resource_name` - instance name.
      * `resource` - instance resource.
      * `self_built` - whether the cluster is built by yourself instead of cloud product.
      * `start_time` - when Offset type timestamp is required.
      * `table_mappings` - maps of table to topic, required when multi topic is selected.
        * `database` - database name.
        * `table` - table name.
        * `topic_id` - Topic ID.
        * `topic` - Topic name.
      * `topic_id` - Topic Id.
      * `topic` - Topic name, use `,` when more than 1 topic.
      * `use_auto_create_topic` - Does the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).
      * `use_table_mapping` - whether to use multi table.
      * `zone_id` - Zone ID.
    * `maria_db_param` - MariaDB configuration, Required when Type is MARIADB.
      * `database` - MariaDB database name, * for all database.
      * `include_content_changes` - If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.
      * `include_query` - If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.
      * `is_table_prefix` - When the Table input is a prefix, the value of this item is true, otherwise it is false.
      * `key_columns` - Format  library 1. table 1: field 1, field 2; library 2. table 2: field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.
      * `output_format` - output format, DEFAULT, CANAL_1, CANAL_2.
      * `record_with_schema` - If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.
      * `resource` - MariaDB connection Id.
      * `snapshot_mode` - schema_only|initial, default initial.
      * `table` - MariaDB db name, is the non-system table in all the monitored databases, you can use to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.
    * `mongo_db_param` - MongoDB config, Required when Type is MONGODB.
      * `collection` - MongoDB collection.
      * `copy_existing` - Whether to copy the stock data, the default parameter is true.
      * `database` - MongoDB database name.
      * `ip` - Mongo DB connection ip.
      * `listening_event` - Listening event type, if it is empty, it means select all. Values include insert, update, replace, delete, invalidate, drop, dropdatabase, rename, used between multiple types, separated by commas.
      * `password` - MongoDB database password.
      * `pipeline` - aggregation pipeline.
      * `port` - MongoDB connection port.
      * `read_preference` - Master-slave priority, default master node.
      * `resource` - resource id.
      * `self_built` - Whether it is a self-built cluster.
      * `user_name` - MongoDB database user name.
    * `my_sql_param` - MySQL configuration, Required when Type is MYSQL.
      * `data_source_increment_column` - The name of the column to be monitored.
      * `data_source_increment_mode` - TIMESTAMP indicates that the incremental column is of timestamp type, INCREMENT indicates that the incremental column is of self-incrementing id type&#39;.
      * `data_source_monitor_mode` - TABLE indicates that the read item is a table, QUERY indicates that the read item is a query.
      * `data_source_monitor_resource` - When DataMonitorMode=TABLE, pass in the Table that needs to be read; when DataMonitorMode=QUERY, pass in the query sql statement that needs to be read.
      * `data_source_start_from` - HEAD means copy stock + incremental data, TAIL means copy only incremental data.
      * `data_target_insert_mode` - INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.
      * `data_target_primary_key_field` - When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.
      * `data_target_record_mapping` - Mapping relationship between tables and messages.
        * `allow_null` - Whether the message is allowed to be empty.
        * `auto_increment` - Whether it is an auto-increment column.
        * `column_name` - Corresponding mapping column name.
        * `column_size` - current column size.
        * `decimal_digits` - current column precision.
        * `default_value` - Database table default parameters.
        * `extra_info` - Database table extra fields.
        * `json_key` - The key name of the message.
        * `type` - message type.
      * `database` - MySQL database name, * is the whole database.
      * `ddl_topic` - The Topic that stores the Ddl information of My SQL, if it is empty, it will not be stored by default.
      * `drop_cls` - When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.
        * `drop_cls_log_set` - cls LogSet id.
        * `drop_cls_owneruin` - account.
        * `drop_cls_region` - The region where the cls is delivered.
        * `drop_cls_topic_id` - cls topic.
        * `drop_invalid_message_to_cls` - Whether to deliver to cls.
      * `drop_invalid_message` - Whether to discard messages that fail to parse, the default is true.
      * `include_content_changes` - If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.
      * `include_query` - If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.
      * `is_table_prefix` - When the Table input is a prefix, the value of this item is true, otherwise it is false.
      * `is_table_regular` - Whether the input table is a regular expression, if this option and Is Table Prefix are true at the same time, the judgment priority of this option is higher than Is Table Prefix.
      * `key_columns` - Format library1.table1 field 1,field 2;library 2.table2 field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.
      * `output_format` - output format, DEFAULT, CANAL_1, CANAL_2.
      * `record_with_schema` - If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.
      * `resource` - MySQL connection Id.
      * `signal_database` - database name of signal table.
      * `snapshot_mode` - whether to Copy inventory information (schema_only does not copy, initial full amount), the default is initial.
      * `table` - The name of the MySQL data table,  is the non-system table in all the monitored databases, which can be separated by, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name, when a regular expression needs to be filled in, the format is data database name.data table name.
      * `topic_regex` - Regular expression for routing events to specific topics, defaults to (.*).
      * `topic_replacement` - TopicRegex, $1, $2.
    * `postgre_sql_param` - PostgreSQL configuration, Required when Type is POSTGRESQL or TDSQL C_POSTGRESQL.
      * `data_format` - Upstream data format (JSON|Debezium), required when the database synchronization mode matches the default field.
      * `data_target_insert_mode` - INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.
      * `data_target_primary_key_field` - When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.
      * `data_target_record_mapping` - Mapping relationship between tables and messages.
        * `allow_null` - Whether the message is allowed to be empty.
        * `auto_increment` - Whether it is an auto-increment column.
        * `column_name` - Column Name.
        * `column_size` - current ColumnSize.
        * `decimal_digits` - current Column DecimalDigits.
        * `default_value` - Database table default parameters.
        * `extra_info` - Database table extra fields.
        * `json_key` - The key name of the message.
        * `type` - message type.
      * `database` - PostgreSQL database name.
      * `drop_invalid_message` - Whether to discard messages that fail to parse, the default is true.
      * `is_table_regular` - Whether the input table is a regular expression.
      * `key_columns` - Format  library1.table1:field 1,field2;library2.table2:field2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.
      * `plugin_name` - (decoderbufs/pgoutput), default decoderbufs.
      * `record_with_schema` - If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.
      * `resource` - PostgreSQL connection Id.
      * `snapshot_mode` - never|initial, default initial.
      * `table` - PostgreSQL tableName, * is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of Schema name.Data table name, and you need to fill in a regular expression When, the format is Schema name.data table name.
    * `scf_param` - Scf configuration, Required when Type is SCF.
      * `batch_size` - The maximum number of messages sent in each batch, the default is 1000.
      * `function_name` - SCF function name.
      * `max_retries` - The number of retries after the SCF call fails, the default is 5.
      * `namespace` - SCF cloud function namespace, the default is default.
      * `qualifier` - SCF cloud function version and alias, the default is DEFAULT.
    * `sql_server_param` - SQLServer configuration, Required when Type is SQLSERVER.
      * `database` - SQLServer database name.
      * `resource` - SQLServer connection Id.
      * `snapshot_mode` - schema_only|initial default initial.
      * `table` - SQLServer table is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.
    * `tdw_param` - Tdw configuration, required when Type is TDW.
      * `bid` - Tdw bid.
      * `is_domestic` - default true.
      * `tdw_host` - TDW address, defalt tl-tdbank-tdmanager.tencent-distribute.com.
      * `tdw_port` - TDW port, default 8099.
      * `tid` - Tdw tid.
    * `topic_param` - Topic configuration, Required when Type is Topic.
      * `compression_type` - Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.
      * `msg_multiple` - 1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).
      * `offset_type` - Offset type, initial position earliest, latest position latest, time point position timestamp.
      * `resource` - The topic name of the topic sold separately.
      * `start_time` - It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.
      * `topic_id` - Topic TopicId.
      * `use_auto_create_topic` - whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).
    * `type` - resource type.
  * `status` - Status, -1 failed to create, 0 to create, 1 to run, 2 to delete, 3 to deleted, 4 to delete failed, 5 to pause, 6 to pause, 7 to pause, 8 to resume, 9 to resume failed.
  * `step_list` - StepList.
  * `target_resource` - Target Resource.
    * `click_house_param` - ClickHouse config, Type CLICKHOUSE requierd.
      * `cluster` - ClickHouse cluster.
      * `database` - ClickHouse database name.
      * `drop_cls` - When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.
        * `drop_cls_log_set` - cls LogSet id.
        * `drop_cls_owneruin` - cls account.
        * `drop_cls_region` - cls region.
        * `drop_cls_topic_id` - cls topicId.
        * `drop_invalid_message_to_cls` - Whether to deliver to cls.
      * `drop_invalid_message` - Whether ClickHouse discards the message that fails to parse, the default is true.
      * `ip` - ClickHouse ip.
      * `password` - ClickHouse passwd.
      * `port` - ClickHouse port.
      * `resource` - resource id.
      * `schema` - ClickHouse schema.
        * `allow_null` - Whether the column item is allowed to be empty.
        * `column_name` - column name.
        * `json_key` - The json Key name corresponding to this column.
        * `type` - type of table column.
      * `self_built` - Whether it is a self-built cluster.
      * `service_vip` - instance vip.
      * `table` - ClickHouse table.
      * `type` - ClickHouse type, emr-clickhouse: emr; cdw-clickhouse: cdwch; selfBuilt: "".
      * `uniq_vpc_id` - instance vpc id.
      * `user_name` - ClickHouse user name.
    * `cls_param` - Cls configuration, Required when Type is CLS.
      * `content_key` - Required when Decode Json is false.
      * `decode_json` - Whether the produced information is in json format.
      * `log_set` - LogSet id.
      * `resource` - cls id.
      * `time_field` - Specify the content of a field in the message as the time of the cls log. The format of the field content needs to be a second-level timestamp.
    * `cos_param` - Cos configuration, required when Type is COS.
      * `aggregate_batch_size` - The size of aggregated messages MB.
      * `aggregate_interval` - time interval.
      * `bucket_name` - cos bucket name.
      * `directory_time_format` - Partition format formatted according to strptime time.
      * `format_output_type` - The file format after message aggregation csv|json.
      * `object_key_prefix` - Dumped object directory prefix.
      * `object_key` - ObjectKey.
      * `region` - region code.
    * `ctsdb_param` - Ctsdb configuration, Required when Type is CTSDB.
      * `ctsdb_metric` - Ctsdb metric.
      * `resource` - resource id.
    * `dts_param` - Dts configuration, required when Type is DTS.
      * `group_id` - Dts consumer group Id.
      * `group_password` - Dts consumer group passwd.
      * `group_user` - Dts account.
      * `ip` - Dts connection ip.
      * `port` - Dts connection port.
      * `resource` - Dts instance Id.
      * `topic` - Dts topic.
      * `tran_sql` - False to synchronize the original data, true to synchronize the parsed json format data, the default is true.
    * `es_param` - Es configuration, required when Type is ES.
      * `content_key` - key for data in non-json format.
      * `database_primary_key` - When the message dumped to ES is the binlog of Database, if you need to synchronize database operations, that is, fill in the primary key of the database table when adding, deleting, and modifying operations to ES.
      * `date_format` - Es date suffix.
      * `document_id_field` - The field name of the document ID value dumped into Es.
      * `drop_cls` - When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.
        * `drop_cls_log_set` - cls log set.
        * `drop_cls_owneruin` - Delivery account of cls.
        * `drop_cls_region` - The region where the cls is delivered.
        * `drop_cls_topic_id` - topic of cls.
        * `drop_invalid_message_to_cls` - Whether to deliver to cls.
      * `drop_dlq` - dead letter queue.
        * `dlq_type` - dlq type, CKAFKA|TOPIC.
        * `kafka_param` - Ckafka type dlq.
          * `compression_type` - Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.
          * `connector_sync_type` - ConnectorSyncType.
          * `enable_toleration` - Enable the fault-tolerant instance and enable the dead-letter queue.
          * `keep_partition` - KeepPartition.
          * `msg_multiple` - 1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).
          * `offset_type` - Offset type, initial position earliest, latest position latest, time point position timestamp.
          * `partition_num` - Partition num.
          * `qps_limit` - Qps limit.
          * `resource_name` - resource id name.
          * `resource` - resource id.
          * `self_built` - Whether it is a self-built cluster.
          * `start_time` - It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.
          * `table_mappings` - The route from Table to Topic must be passed when the Distribute to multiple topics switch is turned on.
            * `database` - database name.
            * `table` - Table name, multiple tables, separated by (commas).
            * `topic_id` - Topic ID.
            * `topic` - Topic name.
          * `topic_id` - Topic Id.
          * `topic` - Topic name, multiple separated by.
          * `use_auto_create_topic` - whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).
          * `use_table_mapping` - Distribute to multiple topics switch, the default is false.
          * `zone_id` - Zone ID.
        * `max_retry_attempts` - retry times.
        * `retry_interval` - retry interval.
        * `topic_param` - DIP Topic type dead letter queue.
          * `compression_type` - Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.
          * `msg_multiple` - 1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).
          * `offset_type` - Offset type, initial position earliest, latest position latest, time point position timestamp.
          * `resource` - The topic name of the topic sold separately.
          * `start_time` - It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.
          * `topic_id` - TopicId.
          * `use_auto_create_topic` - whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).
        * `type` - type, DLQ dead letter queue, IGNORE_ERROR|DROP.
      * `drop_invalid_json_message` - Whether Es discards messages in non-json format.
      * `drop_invalid_message` - Whether Es discards the message of parsing failure.
      * `index_type` - Es custom index name type, STRING, JSONPATH, the default is STRING.
      * `index` - Es index name.
      * `password` - Es Password.
      * `port` - Es connection port.
      * `resource` - Resource.
      * `self_built` - Whether it is a self-built cluster.
      * `service_vip` - instance vip.
      * `uniq_vpc_id` - instance vpc id.
      * `user_name` - Es UserName.
    * `event_bus_param` - EB configuration, required when type is EB.
      * `function_name` - SCF function name.
      * `namespace` - SCF namespace.
      * `qualifier` - SCF version and alias.
      * `resource` - instance id.
      * `self_built` - Whether it is a self-built cluster.
      * `type` - resource type EB_COS/EB_ES/EB_CLS.
    * `kafka_param` - ckafka configuration, required when Type is KAFKA.
      * `compression_type` - Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.
      * `connector_sync_type` - ConnectorSyncType.
      * `enable_toleration` - enable dead letter queue.
      * `keep_partition` - KeepPartition.
      * `msg_multiple` - 1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).
      * `offset_type` - Offset type, from beginning:earliest, from latest:latest, from specific time:timestamp.
      * `partition_num` - the partition num of the topic.
      * `qps_limit` - Qps(query per seconds) limit.
      * `resource_name` - instance name.
      * `resource` - instance resource.
      * `self_built` - whether the cluster is built by yourself instead of cloud product.
      * `start_time` - when Offset type timestamp is required.
      * `table_mappings` - maps of table to topic, required when multi topic is selected.
        * `database` - database name.
        * `table` - table name.
        * `topic_id` - Topic ID.
        * `topic` - Topic name.
      * `topic_id` - Topic Id.
      * `topic` - Topic name, use `,` when more than 1 topic.
      * `use_auto_create_topic` - Does the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).
      * `use_table_mapping` - whether to use multi table.
      * `zone_id` - Zone ID.
    * `maria_db_param` - MariaDB configuration, Required when Type is MARIADB.
      * `database` - MariaDB database name, * for all database.
      * `include_content_changes` - If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.
      * `include_query` - If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.
      * `is_table_prefix` - When the Table input is a prefix, the value of this item is true, otherwise it is false.
      * `key_columns` - Format  library 1. table 1: field 1, field 2; library 2. table 2: field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.
      * `output_format` - output format, DEFAULT, CANAL_1, CANAL_2.
      * `record_with_schema` - If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.
      * `resource` - MariaDB connection Id.
      * `snapshot_mode` - schema_only|initial, default initial.
      * `table` - MariaDB db name, is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.
    * `mongo_db_param` - MongoDB config, Required when Type is MONGODB.
      * `collection` - MongoDB collection.
      * `copy_existing` - Whether to copy the stock data, the default parameter is true.
      * `database` - MongoDB database name.
      * `ip` - Mongo DB connection ip.
      * `listening_event` - Listening event type, if it is empty, it means select all. Values include insert, update, replace, delete, invalidate, drop, dropdatabase, rename, used between multiple types, separated by commas.
      * `password` - MongoDB database password.
      * `pipeline` - aggregation pipeline.
      * `port` - MongoDB connection port.
      * `read_preference` - Master-slave priority, default master node.
      * `resource` - resource id.
      * `self_built` - Whether it is a self-built cluster.
      * `user_name` - MongoDB database user name.
    * `my_sql_param` - MySQL configuration, Required when Type is MYSQL.
      * `data_source_increment_column` - the name of the column to be monitored.
      * `data_source_increment_mode` - TIMESTAMP indicates that the incremental column is of timestamp type, INCREMENT indicates that the incremental column is of self-incrementing id type.
      * `data_source_monitor_mode` - TABLE indicates that the read item is a table, QUERY indicates that the read item is a query.
      * `data_source_monitor_resource` - When DataMonitorMode=TABLE, pass in the Table that needs to be read; when DataMonitorMode=QUERY, pass in the query sql statement that needs to be read.
      * `data_source_start_from` - HEAD means copy stock + incremental data, TAIL means copy only incremental data.
      * `data_target_insert_mode` - INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.
      * `data_target_primary_key_field` - When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.
      * `data_target_record_mapping` - Mapping relationship between tables and messages.
        * `allow_null` - Whether the message is allowed to be empty.
        * `auto_increment` - Whether it is an auto-increment column.
        * `column_name` - Corresponding mapping column name.
        * `column_size` - current column size.
        * `decimal_digits` - current column precision.
        * `default_value` - Database table default parameters.
        * `extra_info` - Database table extra fields.
        * `json_key` - The key name of the message.
        * `type` - message type.
      * `database` - MySQL database name, * is the whole database.
      * `ddl_topic` - The Topic that stores the Ddl information of My SQL, if it is empty, it will not be stored by default.
      * `drop_cls` - When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.
        * `drop_cls_log_set` - cls LogSet id.
        * `drop_cls_owneruin` - account.
        * `drop_cls_region` - The region where the cls is delivered.
        * `drop_cls_topic_id` - cls topic.
        * `drop_invalid_message_to_cls` - Whether to deliver to cls.
      * `drop_invalid_message` - Whether to discard messages that fail to parse, the default is true.
      * `include_content_changes` - If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.
      * `include_query` - If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.
      * `is_table_prefix` - When the Table input is a prefix, the value of this item is true, otherwise it is false.
      * `is_table_regular` - Whether the input table is a regular expression, if this option and Is Table Prefix are true at the same time, the judgment priority of this option is higher than Is Table Prefix.
      * `key_columns` - Format library1.table1 field 1,field 2;library 2.table2 field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.
      * `output_format` - output format, DEFAULT, CANAL_1, CANAL_2.
      * `record_with_schema` - If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.
      * `resource` - MySQL connection Id.
      * `signal_database` - database name of signal table.
      * `snapshot_mode` - whether to Copy inventory information (schema_only does not copy, initial full amount), the default is initial.
      * `table` - The name of the MySQL data table is the non-system table in all the monitored databases, which can be separated by, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name, when a regular expression needs to be filled in, the format is data database name.data table name.
      * `topic_regex` - Regular expression for routing events to specific topics, defaults to (.*).
      * `topic_replacement` - TopicRegex, $1, $2.
    * `postgre_sql_param` - PostgreSQL configuration, Required when Type is POSTGRESQL or TDSQL C_POSTGRESQL.
      * `data_format` - Upstream data format (JSON|Debezium), required when the database synchronization mode matches the default field.
      * `data_target_insert_mode` - INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.
      * `data_target_primary_key_field` - When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.
      * `data_target_record_mapping` - Mapping relationship between tables and messages.
        * `allow_null` - Whether the message is allowed to be empty.
        * `auto_increment` - Whether it is an auto-increment column.
        * `column_name` - Column Name.
        * `column_size` - current ColumnSize.
        * `decimal_digits` - current Column DecimalDigits.
        * `default_value` - Database table default parameters.
        * `extra_info` - Database table extra fields.
        * `json_key` - The key name of the message.
        * `type` - message type.
      * `database` - PostgreSQL database name.
      * `drop_invalid_message` - Whether to discard messages that fail to parse, the default is true.
      * `is_table_regular` - Whether the input table is a regular expression.
      * `key_columns` - Format  library1.table1:field 1,field2;library2.table2:field2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.
      * `plugin_name` - (decoderbufs/pgoutput), default decoderbufs.
      * `record_with_schema` - If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.
      * `resource` - PostgreSQL connection Id.
      * `snapshot_mode` - never|initial, default initial.
      * `table` - PostgreSQL tableName, is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of Schema name.Data table name, and you need to fill in a regular expression When, the format is Schema name.data table name.
    * `scf_param` - Scf configuration, Required when Type is SCF.
      * `batch_size` - The maximum number of messages sent in each batch, the default is 1000.
      * `function_name` - SCF function name.
      * `max_retries` - The number of retries after the SCF call fails, the default is 5.
      * `namespace` - SCF cloud function namespace, the default is default.
      * `qualifier` - SCF cloud function version and alias, the default is DEFAULT.
    * `sql_server_param` - SQLServer configuration, Required when Type is SQLSERVER.
      * `database` - SQLServer database name.
      * `resource` - SQLServer connection Id.
      * `snapshot_mode` - schema_only|initial default initial.
      * `table` - SQLServer table, is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.
    * `tdw_param` - Tdw configuration, required when Type is TDW.
      * `bid` - Tdw bid.
      * `is_domestic` - default true.
      * `tdw_host` - TDW address, defalt tl-tdbank-tdmanager.tencent-distribute.com.
      * `tdw_port` - TDW port, default 8099.
      * `tid` - Tdw tid.
    * `topic_param` - Topic configuration, Required when Type is Topic.
      * `compression_type` - Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.
      * `msg_multiple` - 1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).
      * `offset_type` - Offset type, initial position earliest, latest position latest, time point position timestamp.
      * `resource` - The topic name of the topic sold separately.
      * `start_time` - It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.
      * `topic_id` - Topic TopicId.
      * `use_auto_create_topic` - whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).
    * `type` - Resource Type.
  * `task_current_step` - Task Current Step.
  * `task_id` - task ID.
  * `task_name` - TaskName.
  * `task_progress` - Creation progress percentage.
  * `task_type` - TaskType, SOURCE|SINK.


