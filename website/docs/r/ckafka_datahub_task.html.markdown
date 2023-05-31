---
subcategory: "Cloud Kafka(ckafka)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ckafka_datahub_task"
sidebar_current: "docs-tencentcloud-resource-ckafka_datahub_task"
description: |-
  Provides a resource to create a ckafka datahub_task
---

# tencentcloud_ckafka_datahub_task

Provides a resource to create a ckafka datahub_task

## Example Usage

```hcl
resource "tencentcloud_ckafka_datahub_task" "datahub_task" {
  task_name = "test-task123321"
  task_type = "SOURCE"
  source_resource {
    type = "POSTGRESQL"
    postgre_sql_param {
      database           = "postgres"
      table              = "*"
      resource           = "resource-y9nxnw46"
      plugin_name        = "decoderbufs"
      snapshot_mode      = "never"
      is_table_regular   = false
      key_columns        = ""
      record_with_schema = false
    }
  }
  target_resource {
    type = "TOPIC"
    topic_param {
      compression_type      = "none"
      resource              = "1308726196-keep-topic"
      use_auto_create_topic = false
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `task_name` - (Required, String) name of the task.
* `task_type` - (Required, String, ForceNew) type of the task, SOURCE(data input), SINK(data output).
* `schema_id` - (Optional, String, ForceNew) SchemaId.
* `source_resource` - (Optional, List, ForceNew) data resource.
* `target_resource` - (Optional, List, ForceNew) Target Resource.
* `transform_param` - (Optional, List, ForceNew) Data Processing Rules.
* `transforms_param` - (Optional, List, ForceNew) Data processing rules.

The `analyse_result` object supports the following:

* `key` - (Required, String) KEY.
* `operate` - (Required, String) Operation, DATE system preset - timestamp, CUSTOMIZE customization, MAPPING mapping, JSONPATH.
* `scheme_type` - (Required, String) data type, ORIGINAL, STRING, INT64, FLOAT64, BOOLEAN, MAP, ARRAY.
* `original_value` - (Optional, String) OriginalValue.
* `value_operate` - (Optional, List) VALUE process.
* `value_operates` - (Optional, List) VALUE process chain.
* `value` - (Optional, String) VALUE.

The `analyse_result` object supports the following:

* `key` - (Required, String) key.
* `type` - (Optional, String) Type, DEFAULT default, DATE system default - timestamp, CUSTOMIZE custom, MAPPING mapping.
* `value` - (Optional, String) value.

The `analyse` object supports the following:

* `format` - (Required, String) Parsing format, JSON, DELIMITER delimiter, REGULAR regular extraction, SOURCE processing all results of the upper layer.
* `input_value_type` - (Optional, String) KEY to be processed again - mode.
* `input_value` - (Optional, String) KEY to be processed again - KEY expression.
* `regex` - (Optional, String) delimiter, regular expression.

The `batch_analyse` object supports the following:

* `format` - (Required, String) ONE BY ONE single output, MERGE combined output.

The `click_house_param` object supports the following:

* `cluster` - (Required, String) ClickHouse cluster.
* `database` - (Required, String) ClickHouse database name.
* `resource` - (Required, String) resource id.
* `schema` - (Required, List) ClickHouse schema.
* `table` - (Required, String) ClickHouse table.
* `drop_cls` - (Optional, List) When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.
* `drop_invalid_message` - (Optional, Bool) Whether ClickHouse discards the message that fails to parse, the default is true.
* `ip` - (Optional, String) ClickHouse ip.
* `password` - (Optional, String) ClickHouse passwd.
* `port` - (Optional, Int) ClickHouse port.
* `self_built` - (Optional, Bool) Whether it is a self-built cluster.
* `service_vip` - (Optional, String) instance vip.
* `type` - (Optional, String) ClickHouse type, emr-clickhouse: emr;cdw-clickhouse: cdwch; selfBuilt: ``.
* `uniq_vpc_id` - (Optional, String) instance vpc id.
* `user_name` - (Optional, String) ClickHouse user name.

The `click_house_param` object supports the following:

* `cluster` - (Required, String) ClickHouse cluster.
* `database` - (Required, String) ClickHouse database name.
* `resource` - (Required, String) resource id.
* `schema` - (Required, List) ClickHouse schema.
* `table` - (Required, String) ClickHouse table.
* `drop_cls` - (Optional, List) When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.
* `drop_invalid_message` - (Optional, Bool) Whether ClickHouse discards the message that fails to parse, the default is true.
* `ip` - (Optional, String) ClickHouse ip.
* `password` - (Optional, String) ClickHouse passwd.
* `port` - (Optional, Int) ClickHouse port.
* `self_built` - (Optional, Bool) Whether it is a self-built cluster.
* `service_vip` - (Optional, String) instance vip.
* `type` - (Optional, String) ClickHouse type, emr-clickhouse: emr;cdw-clickhouse: cdwch;selfBuilt: ``.
* `uniq_vpc_id` - (Optional, String) instance vpc id.
* `user_name` - (Optional, String) ClickHouse user name.

The `cls_param` object supports the following:

* `decode_json` - (Required, Bool) Whether the produced information is in json format.
* `resource` - (Required, String) cls id.
* `content_key` - (Optional, String) Required when Decode Json is false.
* `log_set` - (Optional, String) LogSet id.
* `time_field` - (Optional, String) Specify the content of a field in the message as the time of the cls log. The format of the field content needs to be a second-level timestamp.

The `cos_param` object supports the following:

* `bucket_name` - (Required, String) cos bucket name.
* `region` - (Required, String) region code.
* `aggregate_batch_size` - (Optional, Int) The size of aggregated messages MB.
* `aggregate_interval` - (Optional, Int) time interval.
* `directory_time_format` - (Optional, String) Partition format formatted according to strptime time.
* `format_output_type` - (Optional, String) The file format after message aggregation csv|json.
* `object_key_prefix` - (Optional, String) Dumped object directory prefix.
* `object_key` - (Optional, String) ObjectKey.

The `ctsdb_param` object supports the following:

* `ctsdb_metric` - (Optional, String) Ctsdb metric.
* `resource` - (Optional, String) resource id.

The `data_target_record_mapping` object supports the following:

* `allow_null` - (Optional, Bool) Whether the message is allowed to be empty.
* `auto_increment` - (Optional, Bool) Whether it is an auto-increment column.
* `column_name` - (Optional, String) Column Name.
* `column_size` - (Optional, String) current ColumnSize.
* `decimal_digits` - (Optional, String) current Column DecimalDigits.
* `default_value` - (Optional, String) Database table default parameters.
* `extra_info` - (Optional, String) Database table extra fields.
* `json_key` - (Optional, String) The key name of the message.
* `type` - (Optional, String) message type.

The `data_target_record_mapping` object supports the following:

* `allow_null` - (Optional, Bool) Whether the message is allowed to be empty.
* `auto_increment` - (Optional, Bool) Whether it is an auto-increment column.
* `column_name` - (Optional, String) Corresponding mapping column name.
* `column_size` - (Optional, String) current column size.
* `decimal_digits` - (Optional, String) current column precision.
* `default_value` - (Optional, String) Database table default parameters.
* `extra_info` - (Optional, String) Database table extra fields.
* `json_key` - (Optional, String) The key name of the message.
* `type` - (Optional, String) message type.

The `date` object supports the following:

* `format` - (Optional, String) Time format.
* `target_type` - (Optional, String) input type, string|unix.
* `time_zone` - (Optional, String) default GMT+8.

The `drop_cls` object supports the following:

* `drop_cls_log_set` - (Optional, String) cls LogSet id.
* `drop_cls_owneruin` - (Optional, String) account.
* `drop_cls_region` - (Optional, String) The region where the cls is delivered.
* `drop_cls_topic_id` - (Optional, String) cls topic.
* `drop_invalid_message_to_cls` - (Optional, Bool) Whether to deliver to cls.

The `drop_cls` object supports the following:

* `drop_cls_log_set` - (Optional, String) cls LogSet id.
* `drop_cls_owneruin` - (Optional, String) cls account.
* `drop_cls_region` - (Optional, String) cls region.
* `drop_cls_topic_id` - (Optional, String) cls topicId.
* `drop_invalid_message_to_cls` - (Optional, Bool) Whether to deliver to cls.

The `drop_cls` object supports the following:

* `drop_cls_log_set` - (Optional, String) cls log set.
* `drop_cls_owneruin` - (Optional, String) Delivery account of cls.
* `drop_cls_region` - (Optional, String) The region where the cls is delivered.
* `drop_cls_topic_id` - (Optional, String) topic of cls.
* `drop_invalid_message_to_cls` - (Optional, Bool) Whether to deliver to cls.

The `drop_dlq` object supports the following:

* `type` - (Required, String) type, DLQ dead letter queue, IGNORE_ERROR|DROP.
* `dlq_type` - (Optional, String) dlq type, CKAFKA|TOPIC.
* `kafka_param` - (Optional, List) Ckafka type dlq.
* `max_retry_attempts` - (Optional, Int) retry times.
* `retry_interval` - (Optional, Int) retry interval.
* `topic_param` - (Optional, List) DIP Topic type dead letter queue.

The `dts_param` object supports the following:

* `resource` - (Required, String) Dts instance Id.
* `group_id` - (Optional, String) Dts consumer group Id.
* `group_password` - (Optional, String) Dts consumer group passwd.
* `group_user` - (Optional, String) Dts account.
* `ip` - (Optional, String) Dts connection ip.
* `port` - (Optional, Int) Dts connection port.
* `topic` - (Optional, String) Dts topic.
* `tran_sql` - (Optional, Bool) False to synchronize the original data, true to synchronize the parsed json format data, the default is true.

The `es_param` object supports the following:

* `resource` - (Required, String) Resource.
* `content_key` - (Optional, String) key for data in non-json format.
* `database_primary_key` - (Optional, String) When the message dumped to ES is the binlog of Database, if you need to synchronize database operations, that is, fill in the primary key of the database table when adding, deleting, and modifying operations to ES.
* `date_format` - (Optional, String) Es date suffix.
* `document_id_field` - (Optional, String) The field name of the document ID value dumped into Es.
* `drop_cls` - (Optional, List) When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.
* `drop_dlq` - (Optional, List) dead letter queue.
* `drop_invalid_json_message` - (Optional, Bool) Whether Es discards messages in non-json format.
* `drop_invalid_message` - (Optional, Bool) Whether Es discards the message of parsing failure.
* `index_type` - (Optional, String) Es custom index name type, STRING, JSONPATH, the default is STRING.
* `index` - (Optional, String) Es index name.
* `password` - (Optional, String) Es Password.
* `port` - (Optional, Int) Es connection port.
* `self_built` - (Optional, Bool) Whether it is a self-built cluster.
* `service_vip` - (Optional, String) instance vip.
* `uniq_vpc_id` - (Optional, String) instance vpc id.
* `user_name` - (Optional, String) Es UserName.

The `event_bus_param` object supports the following:

* `resource` - (Required, String) instance id.
* `self_built` - (Required, Bool) Whether it is a self-built cluster.
* `type` - (Required, String) resource type, EB_COS/EB_ES/EB_CLS.
* `function_name` - (Optional, String) SCF function name.
* `namespace` - (Optional, String) SCF namespace.
* `qualifier` - (Optional, String) SCF version and alias.

The `event_bus_param` object supports the following:

* `resource` - (Required, String) instance id.
* `self_built` - (Required, Bool) Whether it is a self-built cluster.
* `type` - (Required, String) resource type. EB_COS/EB_ES/EB_CLS.
* `function_name` - (Optional, String) SCF function name.
* `namespace` - (Optional, String) SCF namespace.
* `qualifier` - (Optional, String) SCF version and alias.

The `failure_param` object supports the following:

* `type` - (Required, String) type, DLQ dead letter queue, IGNORE_ERROR|DROP.
* `dlq_type` - (Optional, String) dlq type, CKAFKA|TOPIC.
* `kafka_param` - (Optional, List) Ckafka type dlq.
* `max_retry_attempts` - (Optional, Int) retry times.
* `retry_interval` - (Optional, Int) retry interval.
* `topic_param` - (Optional, List) DIP Topic type dead letter queue.

The `field_chain` object supports the following:

* `analyse` - (Required, List) analyze.
* `analyse_json_result` - (Optional, String) Parsing results in JSON format.
* `analyse_result` - (Optional, List) Analysis result.
* `result` - (Optional, String) Test Results.
* `s_m_t` - (Optional, List) data processing.
* `secondary_analyse_json_result` - (Optional, String) Secondary parsing results in JSON format.
* `secondary_analyse_result` - (Optional, List) Secondary Analysis Results.
* `secondary_analyse` - (Optional, List) secondary analysis.

The `filter_param` object supports the following:

* `key` - (Required, String) Key.
* `match_mode` - (Required, String) Matching mode, prefix matches PREFIX, suffix matches SUFFIX, contains matches CONTAINS, except matches EXCEPT, value matches NUMBER, IP matches IP.
* `value` - (Required, String) Value.
* `type` - (Optional, String) REGULAR.

The `json_path_replace` object supports the following:

* `new_value` - (Required, String) Replacement value, Jsonpath expression or string.
* `old_value` - (Required, String) Replaced value, Jsonpath expression.

The `k_v` object supports the following:

* `delimiter` - (Required, String) delimiter.
* `regex` - (Required, String) Key-value secondary analysis delimiter.
* `keep_original_key` - (Optional, String) Keep the source Key, the default is false not to keep.

The `kafka_param` object supports the following:

* `resource` - (Required, String) instance resource.
* `self_built` - (Required, Bool) whether the cluster is built by yourself instead of cloud product.
* `compression_type` - (Optional, String) Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.
* `enable_toleration` - (Optional, Bool) enable dead letter queue.
* `msg_multiple` - (Optional, Int) 1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).
* `offset_type` - (Optional, String) Offset type, from beginning:earliest, from latest:latest, from specific time:timestamp.
* `partition_num` - (Optional, Int) the partition num of the topic.
* `qps_limit` - (Optional, Int) Qps(query per seconds) limit.
* `resource_name` - (Optional, String) instance name.
* `start_time` - (Optional, Int) when Offset type timestamp is required.
* `table_mappings` - (Optional, List) maps of table to topic, required when multi topic is selected.
* `topic_id` - (Optional, String) Topic id.
* `topic` - (Optional, String) Topic name, use `,` when more than 1 topic.
* `use_auto_create_topic` - (Optional, Bool) Does the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).
* `use_table_mapping` - (Optional, Bool) whether to use multi table.
* `zone_id` - (Optional, Int) Zone ID.

The `kafka_param` object supports the following:

* `resource` - (Required, String) resource id.
* `self_built` - (Required, Bool) Whether it is a self-built cluster.
* `compression_type` - (Optional, String) Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.
* `enable_toleration` - (Optional, Bool) Enable the fault-tolerant instance and enable the dead-letter queue.
* `msg_multiple` - (Optional, Int) 1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).
* `offset_type` - (Optional, String) Offset type, initial position earliest, latest position latest, time point position timestamp.
* `partition_num` - (Optional, Int) Partition num.
* `qps_limit` - (Optional, Int) Qps limit.
* `resource_name` - (Optional, String) resource id name.
* `start_time` - (Optional, Int) It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.
* `table_mappings` - (Optional, List) The route from Table to Topic must be passed when the Distribute to multiple topics switch is turned on.
* `topic_id` - (Optional, String) Topic Id.
* `topic` - (Optional, String) Topic name, multiple separated by `,`.
* `use_auto_create_topic` - (Optional, Bool) whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).
* `use_table_mapping` - (Optional, Bool) Distribute to multiple topics switch, the default is false.
* `zone_id` - (Optional, Int) Zone ID.

The `kafka_param` object supports the following:

* `resource` - (Required, String) resource id.
* `self_built` - (Required, Bool) Whether it is a self-built cluster.
* `compression_type` - (Optional, String) Whether to compress when writing to the Topic, if it is not enabled, fill in none, if it is enabled, fill in open.
* `enable_toleration` - (Optional, Bool) Enable the fault-tolerant instance and enable the dead-letter queue.
* `msg_multiple` - (Optional, Int) 1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).
* `offset_type` - (Optional, String) Offset type, initial position earliest, latest position latest, time point position timestamp.
* `partition_num` - (Optional, Int) Partition num.
* `qps_limit` - (Optional, Int) Qps limit.
* `resource_name` - (Optional, String) resource id name.
* `start_time` - (Optional, Int) It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.
* `table_mappings` - (Optional, List) The route from Table to Topic must be passed when the Distribute to multiple topics switch is turned on.
* `topic_id` - (Optional, String) Topic Id.
* `topic` - (Optional, String) Topic name, multiple separated by,.
* `use_auto_create_topic` - (Optional, Bool) whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks, if you do not use to distribute to multiple topics, you need to fill in the topic name that needs to be automatically created in the Topic field).
* `use_table_mapping` - (Optional, Bool) Distribute to multiple topics switch, the default is false.
* `zone_id` - (Optional, Int) Zone ID.

The `map_param` object supports the following:

* `key` - (Required, String) key.
* `type` - (Optional, String) Type, DEFAULT default, DATE system default - timestamp, CUSTOMIZE custom, MAPPING mapping.
* `value` - (Optional, String) value.

The `maria_db_param` object supports the following:

* `database` - (Required, String) MariaDB database name, * for all database.
* `resource` - (Required, String) MariaDB connection Id.
* `table` - (Required, String) MariaDB db name, *is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.
* `include_content_changes` - (Optional, String) If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.
* `include_query` - (Optional, Bool) If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.
* `is_table_prefix` - (Optional, Bool) When the Table input is a prefix, the value of this item is true, otherwise it is false.
* `key_columns` - (Optional, String) Format  library 1. table 1: field 1, field 2; library 2. table 2: field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.
* `output_format` - (Optional, String) output format, DEFAULT, CANAL_1, CANAL_2.
* `record_with_schema` - (Optional, Bool) If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.
* `snapshot_mode` - (Optional, String) schema_only|initial, default initial.

The `mongo_db_param` object supports the following:

* `collection` - (Required, String) MongoDB collection.
* `copy_existing` - (Required, Bool) Whether to copy the stock data, the default parameter is true.
* `database` - (Required, String) MongoDB database name.
* `resource` - (Required, String) resource id.
* `ip` - (Optional, String) Mongo DB connection ip.
* `listening_event` - (Optional, String) Listening event type, if it is empty, it means select all. Values include insert, update, replace, delete, invalidate, drop, dropdatabase, rename, used between multiple types, separated by commas.
* `password` - (Optional, String) MongoDB database password.
* `pipeline` - (Optional, String) aggregation pipeline.
* `port` - (Optional, Int) MongoDB connection port.
* `read_preference` - (Optional, String) Master-slave priority, default master node.
* `self_built` - (Optional, Bool) Whether it is a self-built cluster.
* `user_name` - (Optional, String) MongoDB database user name.

The `my_sql_param` object supports the following:

* `database` - (Required, String) MySQL database name, * is the whole database.
* `resource` - (Required, String) MySQL connection Id.
* `table` - (Required, String) The name of the MySQL data table,  is the non-system table in all the monitored databases, which can be separated by, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name, when a regular expression needs to be filled in, the format is data database name.data table name.
* `data_source_increment_column` - (Optional, String) The name of the column to be monitored.
* `data_source_increment_mode` - (Optional, String) TIMESTAMP indicates that the incremental column is of timestamp type, INCREMENT indicates that the incremental column is of self-incrementing id type&#39;.
* `data_source_monitor_mode` - (Optional, String) TABLE indicates that the read item is a table, QUERY indicates that the read item is a query.
* `data_source_monitor_resource` - (Optional, String) When DataMonitorMode=TABLE, pass in the Table that needs to be read; when DataMonitorMode=QUERY, pass in the query sql statement that needs to be read.
* `data_source_start_from` - (Optional, String) HEAD means copy stock + incremental data, TAIL means copy only incremental data.
* `data_target_insert_mode` - (Optional, String) INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.
* `data_target_primary_key_field` - (Optional, String) When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.
* `data_target_record_mapping` - (Optional, List) Mapping relationship between tables and messages.
* `ddl_topic` - (Optional, String) The Topic that stores the Ddl information of My SQL, if it is empty, it will not be stored by default.
* `drop_cls` - (Optional, List) When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.
* `drop_invalid_message` - (Optional, Bool) Whether to discard messages that fail to parse, the default is true.
* `include_content_changes` - (Optional, String) If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.
* `include_query` - (Optional, Bool) If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.
* `is_table_prefix` - (Optional, Bool) When the Table input is a prefix, the value of this item is true, otherwise it is false.
* `is_table_regular` - (Optional, Bool) Whether the input table is a regular expression, if this option and Is Table Prefix are true at the same time, the judgment priority of this option is higher than Is Table Prefix.
* `key_columns` - (Optional, String) Format library1.table1 field 1,field 2;library 2.table2 field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.
* `output_format` - (Optional, String) output format, DEFAULT, CANAL_1, CANAL_2.
* `record_with_schema` - (Optional, Bool) If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.
* `signal_database` - (Optional, String) database name of signal table.
* `snapshot_mode` - (Optional, String) whether to Copy inventory information (schema_only does not copy, initial full amount), the default is initial.
* `topic_regex` - (Optional, String) Regular expression for routing events to specific topics, defaults to (.*).
* `topic_replacement` - (Optional, String) TopicRegex, $1, $2.

The `my_sql_param` object supports the following:

* `database` - (Required, String) MySQL database name, * is the whole database.
* `resource` - (Required, String) MySQL connection Id.
* `table` - (Required, String) The name of the MySQL data table,  is the non-system table in all the monitored databases, which can be separated by, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name, when a regular expression needs to be filled in, the format is data database name.data table name.
* `data_source_increment_column` - (Optional, String) the name of the column to be monitored.
* `data_source_increment_mode` - (Optional, String) TIMESTAMP indicates that the incremental column is of timestamp type, INCREMENT indicates that the incremental column is of self-incrementing id type.
* `data_source_monitor_mode` - (Optional, String) TABLE indicates that the read item is a table, QUERY indicates that the read item is a query.
* `data_source_monitor_resource` - (Optional, String) When DataMonitorMode=TABLE, pass in the Table that needs to be read; when DataMonitorMode=QUERY, pass in the query sql statement that needs to be read.
* `data_source_start_from` - (Optional, String) HEAD means copy stock + incremental data, TAIL means copy only incremental data.
* `data_target_insert_mode` - (Optional, String) INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.
* `data_target_primary_key_field` - (Optional, String) When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.
* `data_target_record_mapping` - (Optional, List) Mapping relationship between tables and messages.
* `ddl_topic` - (Optional, String) The Topic that stores the Ddl information of My SQL, if it is empty, it will not be stored by default.
* `drop_cls` - (Optional, List) When the member parameter Drop Invalid Message To Cls is set to true, the Drop Invalid Message parameter is invalid.
* `drop_invalid_message` - (Optional, Bool) Whether to discard messages that fail to parse, the default is true.
* `include_content_changes` - (Optional, String) If the value is all, DDL data and DML data will also be written to the selected topic; if the value is dml, only DML data will be written to the selected topic.
* `include_query` - (Optional, Bool) If the value is true, and the value of the binlog rows query log events configuration item in My SQL is ON, the data flowing into the topic contains the original SQL statement; if the value is false, the data flowing into the topic does not contain Original SQL statement.
* `is_table_prefix` - (Optional, Bool) When the Table input is a prefix, the value of this item is true, otherwise it is false.
* `is_table_regular` - (Optional, Bool) Whether the input table is a regular expression, if this option and Is Table Prefix are true at the same time, the judgment priority of this option is higher than Is Table Prefix.
* `key_columns` - (Optional, String) Format library1.table1 field 1,field 2;library 2.table2 field 2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.
* `output_format` - (Optional, String) output format, DEFAULT, CANAL_1, CANAL_2.
* `record_with_schema` - (Optional, Bool) If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.
* `signal_database` - (Optional, String) database name of signal table.
* `snapshot_mode` - (Optional, String) whether to Copy inventory information (schema_only does not copy, initial full amount), the default is initial.
* `topic_regex` - (Optional, String) Regular expression for routing events to specific topics, defaults to (.*).
* `topic_replacement` - (Optional, String) TopicRegex, $1, $2.

The `postgre_sql_param` object supports the following:

* `database` - (Required, String) PostgreSQL database name.
* `plugin_name` - (Required, String) (decoderbufs/pgoutput), default decoderbufs.
* `resource` - (Required, String) PostgreSQL connection Id.
* `table` - (Required, String) PostgreSQL tableName, * is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of Schema name.Data table name, and you need to fill in a regular expression When, the format is Schema name.data table name.
* `data_format` - (Optional, String) Upstream data format (JSON|Debezium), required when the database synchronization mode matches the default field.
* `data_target_insert_mode` - (Optional, String) INSERT means insert using Insert mode, UPSERT means insert using Upsert mode.
* `data_target_primary_key_field` - (Optional, String) When DataInsertMode=UPSERT, pass in the primary key that the current upsert depends on.
* `data_target_record_mapping` - (Optional, List) Mapping relationship between tables and messages.
* `drop_invalid_message` - (Optional, Bool) Whether to discard messages that fail to parse, the default is true.
* `is_table_regular` - (Optional, Bool) Whether the input table is a regular expression.
* `key_columns` - (Optional, String) Format  library1.table1:field 1,field2;library2.table2:field2, between tables; (semicolon) separated, between fields, (comma) separated. The table that is not specified defaults to the primary key of the table.
* `record_with_schema` - (Optional, Bool) If the value is true, the message will carry the schema corresponding to the message structure, if the value is false, it will not carry.
* `snapshot_mode` - (Optional, String) never|initial, default initial.

The `regex_replace` object supports the following:

* `new_value` - (Required, String) new value.
* `regex` - (Required, String) Regular.

The `replace` object supports the following:

* `new_value` - (Required, String) new value.
* `old_value` - (Required, String) been replaced value.

The `row_param` object supports the following:

* `row_content` - (Required, String) row content, KEY_VALUE, VALUE.
* `entry_delimiter` - (Optional, String) delimiter.
* `key_value_delimiter` - (Optional, String) key, value delimiter.

The `s_m_t` object supports the following:

* `key` - (Required, String) KEY.
* `operate` - (Required, String) Operation, DATE system preset - timestamp, CUSTOMIZE customization, MAPPING mapping, JSONPATH.
* `scheme_type` - (Required, String) data type, ORIGINAL, STRING, INT64, FLOAT64, BOOLEAN, MAP, ARRAY.
* `original_value` - (Optional, String) OriginalValue.
* `value_operate` - (Optional, List) VALUE process.
* `value_operates` - (Optional, List) VALUE process chain.
* `value` - (Optional, String) VALUE.

The `scf_param` object supports the following:

* `function_name` - (Required, String) SCF function name.
* `batch_size` - (Optional, Int) The maximum number of messages sent in each batch, the default is 1000.
* `max_retries` - (Optional, Int) The number of retries after the SCF call fails, the default is 5.
* `namespace` - (Optional, String) SCF cloud function namespace, the default is default.
* `qualifier` - (Optional, String) SCF cloud function version and alias, the default is DEFAULT.

The `schema` object supports the following:

* `allow_null` - (Required, Bool) Whether the column item is allowed to be empty.
* `column_name` - (Required, String) column name.
* `json_key` - (Required, String) The json Key name corresponding to this column.
* `type` - (Required, String) type of table column.

The `secondary_analyse_result` object supports the following:

* `key` - (Required, String) KEY.
* `operate` - (Required, String) Operation, DATE system preset - timestamp, CUSTOMIZE customization, MAPPING mapping, JSONPATH.
* `scheme_type` - (Required, String) data type, ORIGINAL, STRING, INT64, FLOAT64, BOOLEAN, MAP, ARRAY.
* `original_value` - (Optional, String) OriginalValue.
* `value_operate` - (Optional, List) VALUE process.
* `value_operates` - (Optional, List) VALUE process chain.
* `value` - (Optional, String) VALUE.

The `secondary_analyse` object supports the following:

* `regex` - (Required, String) delimiter.

The `source_resource` object supports the following:

* `type` - (Required, String) resource type.
* `click_house_param` - (Optional, List) ClickHouse config, Type CLICKHOUSE requierd.
* `cls_param` - (Optional, List) Cls configuration, Required when Type is CLS.
* `cos_param` - (Optional, List) Cos configuration, required when Type is COS.
* `ctsdb_param` - (Optional, List) Ctsdb configuration, Required when Type is CTSDB.
* `dts_param` - (Optional, List) Dts configuration, required when Type is DTS.
* `es_param` - (Optional, List) Es configuration, required when Type is ES.
* `event_bus_param` - (Optional, List) EB configuration, required when type is EB.
* `kafka_param` - (Optional, List) ckafka configuration, required when Type is KAFKA.
* `maria_db_param` - (Optional, List) MariaDB configuration, Required when Type is MARIADB.
* `mongo_db_param` - (Optional, List) MongoDB config, Required when Type is MONGODB.
* `my_sql_param` - (Optional, List) MySQL configuration, Required when Type is MYSQL.
* `postgre_sql_param` - (Optional, List) PostgreSQL configuration, Required when Type is POSTGRESQL or TDSQL C_POSTGRESQL.
* `scf_param` - (Optional, List) Scf configuration, Required when Type is SCF.
* `sql_server_param` - (Optional, List) SQLServer configuration, Required when Type is SQLSERVER.
* `tdw_param` - (Optional, List) Tdw configuration, required when Type is TDW.
* `topic_param` - (Optional, List) Topic configuration, Required when Type is Topic.

The `split` object supports the following:

* `regex` - (Required, String) delimiter.

The `sql_server_param` object supports the following:

* `database` - (Required, String) SQLServer database name.
* `resource` - (Required, String) SQLServer connection Id.
* `table` - (Required, String) SQLServer table, *is the non-system table in all the monitored databases, you can use, to monitor multiple data tables, but the data table needs to be filled in the format of data database name.data table name.
* `snapshot_mode` - (Optional, String) schema_only|initial default initial.

The `substr` object supports the following:

* `end` - (Required, Int) cut-off position.
* `start` - (Required, Int) interception starting position.

The `table_mappings` object supports the following:

* `database` - (Required, String) database name.
* `table` - (Required, String) Table name, multiple tables, separated by (commas).
* `topic_id` - (Required, String) Topic ID.
* `topic` - (Required, String) Topic name.

The `table_mappings` object supports the following:

* `database` - (Required, String) database name.
* `table` - (Required, String) table name,use, to separate.
* `topic_id` - (Required, String) Topic ID.
* `topic` - (Required, String) Topic name.

The `target_resource` object supports the following:

* `type` - (Required, String) Resource Type.
* `click_house_param` - (Optional, List) ClickHouse config, Type CLICKHOUSE requierd.
* `cls_param` - (Optional, List) Cls configuration, Required when Type is CLS.
* `cos_param` - (Optional, List) Cos configuration, required when Type is COS.
* `ctsdb_param` - (Optional, List) Ctsdb configuration, Required when Type is CTSDB.
* `dts_param` - (Optional, List) Dts configuration, required when Type is DTS.
* `es_param` - (Optional, List) Es configuration, required when Type is ES.
* `event_bus_param` - (Optional, List) EB configuration, required when type is EB.
* `kafka_param` - (Optional, List) ckafka configuration, required when Type is KAFKA.
* `maria_db_param` - (Optional, List) MariaDB configuration, Required when Type is MARIADB.
* `mongo_db_param` - (Optional, List) MongoDB config, Required when Type is MONGODB.
* `my_sql_param` - (Optional, List) MySQL configuration, Required when Type is MYSQL.
* `postgre_sql_param` - (Optional, List) PostgreSQL configuration, Required when Type is POSTGRESQL or TDSQL C_POSTGRESQL.
* `scf_param` - (Optional, List) Scf configuration, Required when Type is SCF.
* `sql_server_param` - (Optional, List) SQLServer configuration, Required when Type is SQLSERVER.
* `tdw_param` - (Optional, List) Tdw configuration, required when Type is TDW.
* `topic_param` - (Optional, List) Topic configuration, Required when Type is Topic.

The `tdw_param` object supports the following:

* `bid` - (Required, String) Tdw bid.
* `tid` - (Required, String) Tdw tid.
* `is_domestic` - (Optional, Bool) default true.
* `tdw_host` - (Optional, String) TDW address, defalt tl-tdbank-tdmanager.tencent-distribute.com.
* `tdw_port` - (Optional, Int) TDW port, default 8099.

The `topic_param` object supports the following:

* `resource` - (Required, String) The topic name of the topic sold separately.
* `compression_type` - (Optional, String) Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.
* `msg_multiple` - (Optional, Int) 1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).
* `offset_type` - (Optional, String) Offset type, initial position earliest, latest position latest, time point position timestamp.
* `start_time` - (Optional, Int) It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.
* `topic_id` - (Optional, String) Topic TopicId.
* `use_auto_create_topic` - (Optional, Bool) whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).

The `topic_param` object supports the following:

* `resource` - (Required, String) The topic name of the topic sold separately.
* `compression_type` - (Optional, String) Whether to perform compression when writing a topic, if it is not enabled, fill in none, if it is enabled, you can choose one of gzip, snappy, lz4 to fill in.
* `msg_multiple` - (Optional, Int) 1 source topic message is amplified into msg Multiple and written to the target topic (this parameter is currently only applicable to ckafka flowing into ckafka).
* `offset_type` - (Optional, String) Offset type, initial position earliest, latest position latest, time point position timestamp.
* `start_time` - (Optional, Int) It must be passed when the Offset type is timestamp, and the time stamp is passed, accurate to the second.
* `topic_id` - (Optional, String) TopicId.
* `use_auto_create_topic` - (Optional, Bool) whether the used topic need to be automatically created (currently only supports SOURCE inflow tasks).

The `transform_param` object supports the following:

* `analysis_format` - (Required, String) parsing format, JSON | DELIMITER| REGULAR.
* `content` - (Required, String) Raw data.
* `failure_param` - (Required, List) Whether to keep parsing failure data.
* `output_format` - (Required, String) output format.
* `source_type` - (Required, String) Data source, TOPIC pulls from the source topic, CUSTOMIZE custom.
* `analyse_result` - (Optional, List) Analysis result.
* `filter_param` - (Optional, List) filter.
* `map_param` - (Optional, List) Map.
* `regex` - (Optional, String) delimiter, regular expression.
* `result` - (Optional, String) Test Results.
* `use_event_bus` - (Optional, Bool) Whether the underlying engine uses eb.

The `transforms_param` object supports the following:

* `content` - (Required, String) Raw data.
* `field_chain` - (Required, List) processing chain.
* `batch_analyse` - (Optional, List) data process.
* `failure_param` - (Optional, List) fail process.
* `filter_param` - (Optional, List) filter.
* `keep_metadata` - (Optional, Bool) Whether to keep the data source Topic metadata information (source Topic, Partition, Offset), the default is false.
* `output_format` - (Optional, String) output format, JSON, ROW, default JSON.
* `result` - (Optional, String) result.
* `row_param` - (Optional, List) The output format is ROW Required.
* `source_type` - (Optional, String) data source.

The `url_decode` object supports the following:

* `charset_name` - (Optional, String) code.

The `value_operate` object supports the following:

* `type` - (Required, String) Processing mode, REPLACE replacement, SUBSTR interception, DATE date conversion, TRIM removal of leading and trailing spaces, REGEX REPLACE regular replacement, URL DECODE, LOWERCASE conversion to lowercase.
* `date` - (Optional, List) Time conversion, required when TYPE=DATE.
* `json_path_replace` - (Optional, List) Json Path replacement, must pass when TYPE=JSON PATH REPLACE.
* `k_v` - (Optional, List) Key-value secondary analysis, must be passed when TYPE=KV.
* `regex_replace` - (Optional, List) Regular replacement, required when TYPE=REGEX REPLACE.
* `replace` - (Optional, List) replace, TYPE=REPLACE is required.
* `result` - (Optional, String) result.
* `split` - (Optional, List) The value supports one split and multiple values, required when TYPE=SPLIT.
* `substr` - (Optional, List) Substr, TYPE=SUBSTR is required.
* `url_decode` - (Optional, List) Url parsing.

The `value_operates` object supports the following:

* `type` - (Required, String) Processing mode, REPLACE replacement, SUBSTR interception, DATE date conversion, TRIM removal of leading and trailing spaces, REGEX REPLACE regular replacement, URL DECODE, LOWERCASE conversion to lowercase.
* `date` - (Optional, List) Time conversion, required when TYPE=DATE.
* `json_path_replace` - (Optional, List) Json Path replacement, must pass when TYPE=JSON PATH REPLACE.
* `k_v` - (Optional, List) Key-value secondary analysis, must be passed when TYPE=KV.
* `regex_replace` - (Optional, List) Regular replacement, required when TYPE=REGEX REPLACE.
* `replace` - (Optional, List) replace, TYPE=REPLACE is required.
* `result` - (Optional, String) result.
* `split` - (Optional, List) The value supports one split and multiple values, required when TYPE=SPLIT.
* `substr` - (Optional, List) Substr, TYPE=SUBSTR is required.
* `url_decode` - (Optional, List) Url parsing.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ckafka datahub_task can be imported using the id, e.g.

```
terraform import tencentcloud_ckafka_datahub_task.datahub_task datahub_task_id
```

