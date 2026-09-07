---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_internal_table"
sidebar_current: "docs-tencentcloud-resource-dlc_internal_table"
description: |-
  Provides a resource to create a DLC internal table
---

# tencentcloud_dlc_internal_table

Provides a resource to create a DLC internal table

## Example Usage

### Create a basic internal table with columns

```hcl
resource "tencentcloud_dlc_internal_table" "example" {
  database_name = "tf_example_database"
  table_name    = "tf_example_table"
  table_comment = "terraform test table"
  type          = "TABLE"
  table_format  = "hive"
  columns {
    name    = "id"
    type    = "bigint"
    comment = "id column"
  }
  columns {
    name    = "name"
    type    = "string"
    comment = "name column"
  }
  properties {
    key   = "format"
    value = "TEXTFILE"
  }
}
```

### Create an Iceberg internal table with smart policy

```hcl
resource "tencentcloud_dlc_internal_table" "example" {
  database_name = "tf_example_database"
  table_name    = "tf_example_iceberg_table"
  table_comment = "terraform iceberg test table"
  type          = "TABLE"
  table_format  = "iceberg"
  primary_keys  = ["id"]
  columns {
    name     = "id"
    type     = "bigint"
    comment  = "id column"
    not_null = true
  }
  columns {
    name    = "name"
    type    = "string"
    comment = "name column"
  }
  smart_policy {
    base_info {
      uin         = "100012345678"
      policy_type = "Table"
      catalog     = "DataLakeCatalog"
      database    = "tf_example_database"
      table       = "tf_example_iceberg_table"
      app_id      = "1300123456"
    }
    policy {
      inherit = "false"
      lifecycle {
        lifecycle_enable     = "enable"
        expiration           = 30
        expired_field        = "dt"
        expired_field_format = "yyyy-MM-dd"
      }
      table_expiration {
        enabled    = true
        expiration = 90
      }
    }
  }
}
```

### Create a table with partitions

```hcl
resource "tencentcloud_dlc_internal_table" "example" {
  database_name = "tf_example_database"
  table_name    = "tf_example_partitioned_table"
  table_comment = "terraform partitioned test table"
  type          = "TABLE"
  table_format  = "hive"
  columns {
    name = "id"
    type = "bigint"
  }
  columns {
    name = "name"
    type = "string"
  }
  partitions {
    name    = "dt"
    type    = "string"
    comment = "date partition"
  }
  properties {
    key   = "format"
    value = "PARQUET"
  }
}
```

## Argument Reference

The following arguments are supported:

* `columns` - (Required, List, ForceNew) Column information of the table.
* `database_name` - (Required, String, ForceNew) The name of the database to which the internal table belongs.
* `table_name` - (Required, String, ForceNew) The name of the internal table.
* `datasource_connection_name` - (Optional, String) The name of the datasource connection to which the table belongs.
* `db_govern_policy_is_disable` - (Optional, String) Whether the database data governance is disabled. true: disabled, false: enabled. Deprecated, use smart_policy instead.
* `govern_policy` - (Optional, List) Data governance policy. Deprecated, use smart_policy instead.
* `partitions` - (Optional, List, ForceNew) Partition information of the table.
* `primary_keys` - (Optional, List: [`String`], ForceNew) Primary keys of the T-ICEBERG table.
* `properties` - (Optional, List, ForceNew) Property information of the table.
* `smart_policy` - (Optional, List) Smart data governance policy.
* `table_comment` - (Optional, String) The comment of the table.
* `table_format` - (Optional, String) The format of the table, such as hive, iceberg, etc.
* `type` - (Optional, String) The type of the table, table or view.
* `upsert_keys` - (Optional, List: [`String`], ForceNew) Upsert keys for V2 upsert table.
* `user_alias` - (Optional, String) The alias of the user who creates the table.
* `user_sub_uin` - (Optional, String) The sub UIN of the user who creates the table.

The `advance_policy` object of `written` supports the following:

* `before_days` - (Optional, Int) Snapshot expiration time in days.
* `compact_enable` - (Optional, String) Whether to enable compaction.
* `compact_strategy` - (Optional, String) File compaction strategy.
* `cow_compact_enable` - (Optional, String) Whether to enable COW table compaction.
* `delete_enable` - (Optional, String) Whether to enable history data cleanup.
* `expired_snapshots_interval_min` - (Optional, Int) Snapshot expiration execution interval in minutes.
* `min_input_files` - (Optional, Int) The minimum number of files to trigger compaction.
* `remove_orphan_interval_min` - (Optional, Int) Orphan file removal execution interval in minutes.
* `retain_last` - (Optional, Int) Number of snapshots to retain.
* `sort_orders` - (Optional, List) Sort compaction strategy rules.
* `target_file_size_bytes` - (Optional, Int) Target file size in bytes for compaction.

The `base_info` object of `smart_policy` supports the following:

* `uin` - (Required, String) User UIN.
* `app_id` - (Optional, String) User APP ID.
* `catalog` - (Optional, String) Catalog name.
* `database` - (Optional, String) Database name.
* `policy_type` - (Optional, String) Policy type: Catalog/Database/Table.
* `table` - (Optional, String) Table name.

The `change_table` object of `policy` supports the following:

* `data_retention_time` - (Optional, Int) Data retention time in days for the change table.

The `columns` object supports the following:

* `name` - (Required, String) Column name.
* `type` - (Required, String) Column type.
* `comment` - (Optional, String) Column comment.
* `default` - (Optional, String) Column default value.
* `is_partition` - (Optional, Bool) Whether the column is a partition field.
* `not_null` - (Optional, Bool) Whether the column is not null.
* `position` - (Optional, Int) Column position, smaller is earlier.
* `precision` - (Optional, Int) The precision of the numeric type.
* `scale` - (Optional, Int) The scale of the numeric type.

The `data_mask_strategy_info` object of `columns` supports the following:


The `favor` object of `resources` supports the following:

* `catalog` - (Optional, String) Catalog name.
* `data_base` - (Optional, String) Database name.
* `priority` - (Optional, Int) Priority.
* `table` - (Optional, String) Table name.

The `govern_policy` object supports the following:

* `govern_engine` - (Optional, String) Governance engine.
* `rule_type` - (Optional, String) Governance rule type. Customize: custom; Intelligence: intelligent governance.

The `index` object of `policy` supports the following:

* `index_enable` - (Optional, String) Whether to enable index.

The `lifecycle` object of `policy` supports the following:

* `drop_table` - (Optional, Bool) Whether to drop the table. Deprecated, use table_expiration instead.
* `expiration` - (Optional, Int) Expiration time.
* `expired_field_format` - (Optional, String) Expired field format.
* `expired_field` - (Optional, String) Expired field.
* `lifecycle_enable` - (Optional, String) Whether lifecycle is enabled.

The `partitions` object supports the following:

* `comment` - (Optional, String) Partition comment.
* `name` - (Optional, String) Partition column name.
* `transform_args` - (Optional, List) Transform strategy arguments.
* `transform` - (Optional, String) Implicit partition transform strategy.
* `type` - (Optional, String) Partition type.

The `policy` object of `smart_policy` supports the following:

* `change_table` - (Optional, List) Change table policy.
* `index` - (Optional, List) Index policy.
* `inherit` - (Optional, String) Whether to inherit.
* `lifecycle` - (Optional, List) Data lifecycle policy.
* `resources` - (Optional, List) Data governance resources.
* `table_expiration` - (Optional, List) Table expiration policy.
* `written` - (Optional, List) Data rewrite policy.

The `properties` object supports the following:

* `key` - (Required, String) Property key.
* `value` - (Required, String) Property value.

The `resource_conf` object of `resources` supports the following:

* `parallelism` - (Optional, Int) The parallelism of the optimization task.

The `resources` object of `policy` supports the following:

* `attribution_type` - (Optional, String) Attribution type.
* `favor` - (Optional, List) Affinity info list.
* `instance` - (Optional, String) Instance name.
* `name` - (Optional, String) Engine name.
* `resource_conf` - (Optional, List) Resource configuration.
* `resource_group_name` - (Optional, String) Standard engine resource group name.
* `resource_type` - (Optional, String) Resource type.
* `status` - (Optional, Int) Status.

The `smart_policy` object supports the following:

* `base_info` - (Optional, List) Base info of the smart policy.
* `policy` - (Optional, List) Smart optimizer policy.

The `sort_orders` object of `advance_policy` supports the following:

* `column` - (Optional, String) Column name for sorting.
* `null_order` - (Optional, String) Null order: at the beginning or end.
* `sort_direction` - (Optional, String) Sort direction: ascending or descending.

The `table_expiration` object of `policy` supports the following:

* `enabled` - (Required, Bool) Whether the policy is enabled.
* `expiration` - (Required, Int) Table expiration time in days.

The `written` object of `policy` supports the following:

* `advance_policy` - (Optional, List) Advanced rewrite policy.
* `written_enable` - (Optional, String) none/enable/disable/default.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Table creation time, in ms.
* `execution` - The auto-generated SQL statement from create.
* `heat_value` - Access heat value.
* `input_format_short` - Abbreviation of InputFormat.
* `input_format` - Data format.
* `is_t_iceberg_sql` - Whether it is T-ICEBERG SQL.
* `location` - Data storage path.
* `map_materialized_view_name` - Materialized view name mapping.
* `modified_time` - Table modification time, in ms.
* `record_count` - Table row count.
* `storage_size` - Table storage size in bytes.

The `columns` object exports the following:

* `create_time` - Column creation time.
* `data_mask_strategy_info` - Data mask strategy info.
  * `strategy_desc` - Strategy description.
  * `strategy_id` - Strategy ID.
  * `strategy_name` - Strategy name.
  * `strategy_type` - Strategy type.
  * `users` - User sub UIN list separated by semicolons.
* `modified_time` - Column modification time.
* `nullable` - Whether the column is nullable.
* `type_text` - Data field description.

The `partitions` object exports the following:

* `create_time` - Partition creation time.


## Import

DLC internal table can be imported using the composite id (`database_name#table_name`), e.g.

```
terraform import tencentcloud_dlc_internal_table.example tf_example_database#tf_example_table
```

