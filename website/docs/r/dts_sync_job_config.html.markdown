---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_sync_job_config"
sidebar_current: "docs-tencentcloud-resource-dts_sync_job_config"
description: |-
  Provides a resource to create a DTS sync job config
---

# tencentcloud_dts_sync_job_config

Provides a resource to create a DTS sync job config

~> **NOTE:** Only `Running` sync jobs can be configured.

## Example Usage

```hcl
resource "tencentcloud_dts_sync_job_config" "example" {
  job_id = "sync-6zm6nj70"

  dynamic_options {
    op_types = ["Insert", "Update", "Delete", "PartialDDL"]
    ddl_options {
      ddl_object = "Database"
      ddl_value  = ["Create", "Drop", "Alter"]
    }

    ddl_options {
      ddl_object = "Table"
      ddl_value  = ["Alter", "Create", "Drop", "Rename", "Truncate"]
    }
    ddl_options {
      ddl_object = "View"
      ddl_value  = ["Create", "Drop"]
    }
    ddl_options {
      ddl_object = "Index"
      ddl_value  = ["Create", "Drop"]
    }

    kafka_option {
      data_option {
        key   = "skipLargeMessage"
        value = "off"
      }
    }

    conflict_handle_type     = "ReportError"
    deal_of_exist_same_table = "ReportErrorAfterCheck"
  }

  dynamic_objects {
    mode = "Partial"
    databases {
      db_name        = "db1"
      db_mode        = "All"
      table_mode     = "All"
      view_mode      = "All"
      procedure_mode = "All"
      function_mode  = "All"
    }

    advanced_objects = ["procedure", "function"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Required, String, ForceNew) Synchronization instance id (i.e. identifies a synchronization job).
* `dynamic_objects` - (Optional, List) Modified sync objects, mirrors SDK Objects struct.
* `dynamic_options` - (Optional, List) Modified sync task options, mirrors SDK DynamicOptions struct.

The `columns` object of `tables` supports the following:

* `column_name` - (Optional, String) Column name. Note: This field may return null, indicating that no valid value can be obtained.
* `new_column_name` - (Optional, String) New column name. Note: This field may return null, indicating that no valid value can be obtained.

The `conflict_handle_option` object of `dynamic_options` supports the following:

* `condition_column` - (Optional, String) Columns covered by the condition. Note: This field may return null, indicating that no valid value can be obtained.
* `condition_operator` - (Optional, String) Conditional Override Operation. Note: This field may return null, indicating that no valid value can be obtained.
* `condition_order_in_src_and_dst` - (Optional, String) Conditional Override Priority Processing. Note: This field may return null, indicating that no valid value can be obtained.

The `data_option` object of `kafka_option` supports the following:

* `key` - (Optional, String) Option key.
* `value` - (Optional, String) Option value.

The `databases_op_filter` object of `dynamic_objects` supports the following:

* `db_name` - (Required, String) Database name that the rule takes effect on.
* `op_filter` - (Optional, List) Database level DDL/DML filter rule.
* `tables` - (Optional, List) Table level DDL/DML filter info.
* `views` - (Optional, List) View level DDL/DML filter info.

The `databases` object of `dynamic_objects` supports the following:

* `db_mode` - (Optional, String) DB selection mode: All (for all objects under the current object), Partial (for some objects), when the Mode is Partial, this item is required. Note that synchronization of advanced objects does not depend on this value. Note: This field may return null, indicating that no valid value can be obtained.
* `db_name` - (Optional, String) The name of the library that needs to be migrated or synchronized. This item is required when the ObjectMode is Partial. Note: This field may return null, indicating that no valid value can be obtained.
* `event_mode` - (Optional, String) Event migration mode, all (for all objects under the current object), partial (partial objects). Note: This field may return null, indicating that no valid value can be obtained.
* `events` - (Optional, Set) When EventMode is partial, specify the name of the event to be migrated. Note: This field may return null, indicating that no valid value can be obtained.
* `function_mode` - (Optional, String) Select the mode to be synchronized, Partial is a part, all is an entire selection. Note: This field may return null, indicating that no valid value can be obtained.
* `functions` - (Optional, Set) Required when the FunctionMode value is Partial. Note: This field may return null, indicating that no valid value can be obtained.
* `new_db_name` - (Optional, String) The name of the library after migration or synchronization, which is the same as the source library by default. Note: This field may return null, indicating that no valid value can be obtained.
* `new_schema_name` - (Optional, String) Schema name after migration or synchronization. Note: This field may return null, indicating that no valid value can be obtained.
* `procedure_mode` - (Optional, String) Select the mode to be synchronized, Partial is part, All is the whole selection. Note: This field may return null, indicating that no valid value can be obtained.
* `procedures` - (Optional, Set) Required when the value of ProcedureMode is Partial. Note: This field may return null, indicating that no valid value can be obtained.
* `schema_mode` - (Optional, String) Schema selection mode, used by PostgreSQL and SQL Server sync links. Valid values: `All` (all objects under the current object), `Partial` (some objects). Note: This field may return null, indicating that no valid value can be obtained.
* `schema_name` - (Optional, String) Migrated or synchronized schema. Note: This field may return null, indicating that no valid value can be obtained.
* `table_mode` - (Optional, String) Table selection mode: All (for all objects under the current object), Partial (for some objects), this item is required when the DBMode is Partial. Note: This field may return null, indicating that no valid value can be obtained.
* `tables` - (Optional, List) A collection of table graph objects, when TableMode is Partial, this item needs to be filled in. Note: This field may return null, indicating that no valid value can be obtained.
* `trigger_mode` - (Optional, String) Trigger migration mode, all (for all objects under the current object), partial (partial objects). Note: This field may return null, indicating that no valid value can be obtained.
* `triggers` - (Optional, Set) When TriggerMode is partial, specify the name of the trigger to be migrated. Note: This field may return null, indicating that no valid value can be obtained.
* `view_mode` - (Optional, String) View selection mode: All is all view objects under the current object, Partial is part of the view objects. Note: This field may return null, indicating that no valid value can be obtained.
* `views` - (Optional, List) View object collection, when ViewMode is Partial, this item needs to be filled in. Note: This field may return null, indicating that no valid value can be obtained.

The `ddl_options` object of `dynamic_options` supports the following:

* `ddl_object` - (Optional, String) Ddl type, such as Database, Table, View, Index, etc. Note: This field may return null, indicating that no valid value can be obtained.
* `ddl_value` - (Optional, Set) The specific value of ddl, the possible values for Database [Create,Drop,Alter].The possible values for Table [Create,Drop,Alter,Truncate,Rename].The possible values for View[Create,Drop].For the possible values of Index [Create, Drop]. Note: This field may return null, indicating that no valid value can be obtained.

The `ddl_options` object of `op_filter` supports the following:

* `ddl_object` - (Optional, String) Ddl type, such as Database, Table, View, Index, etc. Note: This field may return null, indicating that no valid value can be obtained.
* `ddl_value` - (Optional, Set) The specific value of ddl, the possible values for Database [Create,Drop,Alter].The possible values for Table [Create,Drop,Alter,Truncate,Rename].The possible values for View[Create,Drop].For the possible values of Index [Create, Drop]. Note: This field may return null, indicating that no valid value can be obtained.

The `dynamic_objects` object supports the following:

* `advanced_objects` - (Optional, Set) For advanced object types, such as function and procedure, when an advanced object needs to be synchronized, the initialization type must include the structure initialization type, that is, the value of the Options.InitType field is Structure or Full. Note: This field may return null, indicating that no valid value can be obtained.
* `databases_op_filter` - (Optional, List) Library/table/view level DML/DDL white list. Note: This field may return null, indicating that no valid value can be obtained.
* `databases` - (Optional, List) Synchronization object, not null when Mode is Partial. Note: This field may return null, indicating that no valid value can be obtained.
* `mode` - (Optional, String) Migration object type Partial (partial object). Note: This field may return null, indicating that no valid value can be obtained.
* `online_ddl` - (Optional, List, **Deprecated**) It has been deprecated. OnlineDDL type. Note: This field may return null, indicating that no valid value can be obtained.

The `dynamic_options` object supports the following:

* `ddl_options` - (Required, List) DDL synchronization options, specifically describe which DDLs to synchronize. Required field.
* `op_types` - (Required, Set) DML and DDL options to be synchronized, Insert (insert operation), Update (update operation), Delete (delete operation), DDL (structure synchronization), PartialDDL (custom, work with DdlOptions). Required field.
* `conflict_handle_option` - (Optional, List) Detailed options for conflict handling, such as conditional rows and conditional actions in conditional overrides. Note: This field may return null, indicating that no valid value can be obtained.
* `conflict_handle_type` - (Optional, String) Conflict handling options, ReportError (error report, the default value), Ignore (ignore), Cover (cover), ConditionCover (condition coverage). Note: This field may return null, indicating that no valid value can be obtained.
* `deal_of_exist_same_table` - (Optional, String) The processing of the table with the same name, ReportErrorAfterCheck (pre-check and report error, default), ExecuteAfterIgnore (ignore and continue to execute). Note: This field may return null, indicating that no valid value can be obtained.
* `filter_begin_commit` - (Optional, Bool) Whether to filter out begin and commit messages for sync to kafka link. Currently only mysql2kafka link supports it.
* `filter_checkpoint` - (Optional, Bool) Whether to filter out checkpoint messages for sync to kafka link. Currently only mysql2kafka link supports it.
* `kafka_option` - (Optional, List) Kafka config for sync to kafka link. Note: This field may return null, indicating that no valid value can be obtained.
* `start_position` - (Optional, String) Only incremental tasks reset the designated site.

The `kafka_option` object of `dynamic_options` supports the following:

* `data_option` - (Optional, List) Other additional information, for specific data types extra parameters can be set.
* `data_type` - (Optional, String) Data type delivered to kafka, such as Avro, Json, canal-pb, canal-json, debezium.
* `ddl_topic_name` - (Optional, String) Topic used to store ddl.
* `topic_rules` - (Optional, List) Description of single topic and custom topic.
* `topic_type` - (Optional, String) Sync topic strategy, such as Single (centralized delivery to single topic), Multi (custom topic name).

The `online_ddl` object of `dynamic_objects` supports the following:

* `status` - (Required, String, **Deprecated**) It has been deprecated. Status: ON - Enabled, OFF - Disabled.

The `op_filter` object of `databases_op_filter` supports the following:

* `ddl_options` - (Optional, List) DDL white list.
* `op_types` - (Optional, Set) DML white list.

The `op_filter` object of `tables` supports the following:

* `ddl_options` - (Optional, List) DDL white list.
* `op_types` - (Optional, Set) DML white list.

The `op_filter` object of `views` supports the following:

* `ddl_options` - (Optional, List) DDL white list.
* `op_types` - (Optional, Set) DML white list.

The `tables` object of `databases_op_filter` supports the following:

* `table_name` - (Required, String) Table name that the rule takes effect on.
* `op_filter` - (Optional, List) Table level DDL/DML filter rule.

The `tables` object of `databases` supports the following:

* `column_mode` - (Optional, String) Whether to synchronize all columns in the table, All: all columns under the current table, Partial: part of the columns under the current table, fill in the detailed column information through the columns field. Note: This field may return null, indicating that no valid value can be obtained.
* `columns` - (Optional, List) Synchronized column information, required when ColumnMode is Partial. Note: This field may return null, indicating that no valid value can be obtained.
* `filter_condition` - (Optional, String) Filter condition. Note: This field may return null, indicating that no valid value can be obtained.
* `new_table_name` - (Optional, String) New table name. Note: This field may return null, indicating that no valid value can be obtained.
* `table_edit_mode` - (Optional, String) Edit table type, rename (table mapping), pt (synchronize attached tables). Note: This field may return null, indicating that no valid value can be obtained.
* `table_name` - (Optional, String) Table name. Note: This field may return null, indicating that no valid value can be obtained.
* `tmp_tables` - (Optional, Set) Temporary table generated by pt-osc/gh-ost. Note: This field may return null, indicating that no valid value can be obtained.

The `topic_rules` object of `kafka_option` supports the following:

* `columns` - (Optional, Set) Column names for partition by columns, multiple columns can be selected.
* `db_match_mode` - (Optional, String) Database name matching rule, such as Regular, Default.
* `db_name` - (Optional, String) Database name, takes effect when DbMatchMode=Regular.
* `partition_type` - (Optional, String) Topic partition strategy.
* `table_match_mode` - (Optional, String) Table name matching rule, such as Regular, Default.
* `table_name` - (Optional, String) Table name, only takes effect when TableMatchMode=Regular.
* `topic_name` - (Optional, String) Topic name. When single topic, all TopicName must be the same.

The `views` object of `databases_op_filter` supports the following:

* `view_name` - (Required, String) View name that the rule takes effect on.
* `op_filter` - (Optional, List) View level DDL/DML filter rule.

The `views` object of `databases` supports the following:

* `new_view_name` - (Optional, String) New view name. Note: This field may return null, indicating that no valid value can be obtained.
* `view_name` - (Optional, String) View name. Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `actions` - Actions that can be performed in the current state.
* `all_actions` - All supported actions.
* `auto_renew` - Auto renewal flag, 1 (auto renewal), 0 (no auto renewal).
* `auto_retry_time_range_minutes` - Auto retry time period setting.
* `create_time` - Creation time.
* `dst_access_type` - Target access type, cdb (cloud database), cvm (cloud host self-built), vpc (private network), etc.
* `dst_database_type` - Target database type, such as mysql, tdsqlmysql, mariadb, cynosdbmysql, etc.
* `dst_node_type` - Target node type, cluster or single.
* `dst_region` - Target region, such as: ap-guangzhou.
* `dumper_resume_ctrl` - Full export reentrant flag: enum yes/no.
* `end_time` - End time.
* `expect_run_time` - Expected start time, when the value of RunMode is Timed, this value is required, such as: 2006-01-02 15:04:05.
* `expire_time` - Expiration time, format is yyyy-mm-dd hh:mm:ss.
* `instance_class` - Sync link specification, such as micro, small, medium, large.
* `job_name` - Sync job name.
* `offline_time` - Offline time, format is yyyy-mm-dd hh:mm:ss.
* `opt_obj_status` - Dynamic modification object status.
* `pay_mode` - Payment method, PostPay (pay as you go), PrePay (monthly subscription).
* `run_mode` - Operation mode, such as: Immediate (indicates immediate operation, the default value is this value), Timed (indicates scheduled operation).
* `specification` - Task specification.
* `src_access_type` - Source access type, cdb (cloud database), cvm (cloud host self-built), vpc (private network), etc.
* `src_database_type` - Source database type, such as mysql, tdsqlmysql, mariadb, cynosdbmysql, etc.
* `src_node_type` - Source node type, cluster or single.
* `src_region` - Source region, such as: ap-guangzhou.
* `start_time` - Start time.
* `status` - Task status.
* `trade_status` - Status used for billing.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [DescribeSyncJobs](https://www.tencentcloud.com/document/api/571/51793)
* Read: [DescribeSyncJobs](https://www.tencentcloud.com/document/api/571/51793)
* Update:
  - [DescribeSyncJobs](https://www.tencentcloud.com/document/api/571/51793)
  - [ModifySyncJobConfig](https://www.tencentcloud.com/document/api/571/55496)

## Import

DTS sync job config can be imported using the id, e.g.

```
terraform import tencentcloud_dts_sync_job_config.example sync-6zm6nj70
```

