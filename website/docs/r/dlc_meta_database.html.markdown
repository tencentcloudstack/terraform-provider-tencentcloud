---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_meta_database"
sidebar_current: "docs-tencentcloud-resource-dlc_meta_database"
description: |-
  Provides a resource to create a DLC meta database
---

# tencentcloud_dlc_meta_database

Provides a resource to create a DLC meta database

## Example Usage

```hcl
resource "tencentcloud_dlc_meta_database" "example" {
  database_name = "tf_example_db"
  comment       = "tf example meta database"
  govern_policy {
    rule_type     = "Customize"
    govern_engine = "engine_name"
  }
}
```

## Argument Reference

The following arguments are supported:

* `database_name` - (Required, String, ForceNew) Name of the DLC meta database.
* `comment` - (Optional, String) Description of the DLC meta database, length 0~2048.
* `datasource_connection_name` - (Optional, String) Datasource connection name, default `DataLakeCatalog`.
* `govern_policy` - (Optional, List) Data governance config.
* `smart_policy` - (Optional, List) Smart data governance config.

The `advance_policy` object of `written` supports the following:

* `before_days` - (Optional, Int) Snapshot expire before days.
* `compact_enable` - (Optional, String) Whether to enable compact.
* `compact_strategy` - (Optional, String) File compact strategy.
* `cow_compact_enable` - (Optional, String) Whether to enable COW table compact.
* `delete_enable` - (Optional, String) Whether to enable history data cleanup.
* `expired_snapshots_interval_min` - (Optional, Int) Snapshot expire interval in minutes.
* `min_input_files` - (Optional, Int) Min input files to compact.
* `remove_orphan_interval_min` - (Optional, Int) Remove orphan files interval in minutes.
* `retain_last` - (Optional, Int) Retain last snapshot count.
* `sort_orders` - (Optional, List) Sort compact strategy rules.
* `target_file_size_bytes` - (Optional, Int) Target file size after compact.

The `base_info` object of `smart_policy` supports the following:

* `uin` - (Required, String) User UIN.
* `app_id` - (Optional, String) User AppID.
* `catalog` - (Optional, String) Catalog name.
* `database` - (Optional, String) Database name.
* `policy_type` - (Optional, String) Policy type.
* `table` - (Optional, String) Table name.

The `change_table` object of `policy` supports the following:

* `data_retention_time` - (Optional, Int) Data retention time of change table, in days.

The `favor` object of `resources` supports the following:

* `catalog` - (Optional, String) Catalog name.
* `data_base` - (Optional, String) Database name.
* `priority` - (Optional, Int) Priority.
* `table` - (Optional, String) Table name.

The `govern_policy` object supports the following:

* `govern_engine` - (Optional, String) Govern engine.
* `rule_type` - (Optional, String) Govern rule type, `Customize` or `Intelligence`.

The `index` object of `policy` supports the following:

* `index_enable` - (Optional, String) Whether to enable index.

The `lifecycle` object of `policy` supports the following:

* `drop_table` - (Optional, Bool) Whether to drop table (deprecated, use table_expiration instead).
* `expiration` - (Optional, Int) Expiration time.
* `expired_field_format` - (Optional, String) Expired field format.
* `expired_field` - (Optional, String) Expired field.
* `lifecycle_enable` - (Optional, String) Whether to enable lifecycle.

The `policy` object of `smart_policy` supports the following:

* `change_table` - (Optional, List) Change table policy.
* `index` - (Optional, List) Index policy.
* `inherit` - (Optional, String) Whether to inherit.
* `lifecycle` - (Optional, List) Data lifecycle policy.
* `resources` - (Optional, List) Data governance resources.
* `table_expiration` - (Optional, List) Table expiration policy.
* `written` - (Optional, List) Data rewrite policy.

The `resource_conf` object of `resources` supports the following:

* `parallelism` - (Optional, Int) Parallelism of the optimize task.

The `resources` object of `policy` supports the following:

* `attribution_type` - (Optional, String) Attribution type.
* `favor` - (Optional, List) Affinity info.
* `instance` - (Optional, String) Instance name.
* `name` - (Optional, String) Engine name.
* `resource_conf` - (Optional, List) Resource config info.
* `resource_group_name` - (Optional, String) Standard engine resource group name.
* `resource_type` - (Optional, String) Resource type.
* `status` - (Optional, Int) Status.

The `smart_policy` object supports the following:

* `base_info` - (Optional, List) Base info of smart policy.
* `policy` - (Optional, List) Policy description of smart policy.

The `sort_orders` object of `advance_policy` supports the following:

* `column` - (Optional, String) Sort column name.
* `null_order` - (Optional, String) Null order, `first` or `last`.
* `sort_direction` - (Optional, String) Sort direction, `asc` or `desc`.

The `table_expiration` object of `policy` supports the following:

* `enabled` - (Required, Bool) Whether to enable the policy.
* `expiration` - (Required, Int) Table expiration time, in days.

The `written` object of `policy` supports the following:

* `advance_policy` - (Optional, List) User custom advance params.
* `written_enable` - (Optional, String) Whether to enable, `none`/`enable`/`disable`/`default`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `batch_id` - Batch ID of the async task.
* `catalog_name` - Catalog name.
* `catalog_type` - Catalog type.
* `create_time` - Database create time, in seconds.
* `database_id` - Database ID.
* `is_information_schema` - Whether it is InformationSchema.
* `location` - COS storage path.
* `modified_time` - Database modified time, in seconds.
* `properties` - Properties of the database.
  * `key` - Property key.
  * `value` - Property value.
* `task_id_set` - Task ID set.
* `user_alias` - User alias who created the database.
* `user_sub_uin` - User sub UIN who created the database.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `20m`) Used when creating the resource.
* `read` - (Defaults to `10m`) Used when reading the resource.
* `update` - (Defaults to `20m`) Used when updating the resource.
* `delete` - (Defaults to `20m`) Used when deleting the resource.

## Import

DLC meta database can be imported using the datasourceConnectionName#databaseName when `datasource_connection_name` is set, or just the bare databaseName otherwise, e.g.

```
terraform import tencentcloud_dlc_meta_database.example DataLakeCatalog#tf_example_db
```

or

```
terraform import tencentcloud_dlc_meta_database.example tf_example_db
```

