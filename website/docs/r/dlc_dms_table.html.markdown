---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_dms_table"
sidebar_current: "docs-tencentcloud-resource-dlc_dms_table"
description: |-
  Provides a resource to create a DLC DMS table
---

# tencentcloud_dlc_dms_table

Provides a resource to create a DLC DMS table

## Example Usage

```hcl
resource "tencentcloud_dlc_dms_table" "example" {
  db_name = "tf_example_db"
  name    = "tf_example_table"
  type    = "EXTERNAL_TABLE"

  asset {
    name        = "tf_example_table"
    description = "tf example dlc dms table"
    owner       = "root"
  }

  columns {
    name     = "id"
    type     = "bigint"
    position = 1
  }

  columns {
    name     = "name"
    type     = "string"
    position = 2
  }

  sds {
    location      = "cosn://tf-example-bucket/example/"
    input_format  = "org.apache.hadoop.hive.ql.io.avro.AvroContainerInputFormat"
    output_format = "org.apache.hadoop.hive.ql.io.avro.AvroContainerOutputFormat"
    serde_lib     = "org.apache.hadoop.hive.serde2.avro.AvroSerDe"
    serde_params {
      key   = "serialization.format"
      value = "1"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `db_name` - (Required, String) Database name.
* `name` - (Required, String) Table name.
* `asset` - (Optional, List) Basic object.
* `columns` - (Optional, List) Columns.
* `data_update_time` - (Optional, String) Data update time.
* `datasource_connection_name` - (Optional, String) Data source connection name.
* `delete_data` - (Optional, Bool) Whether to delete data.
* `env_props` - (Optional, List) Environment attributes.
* `last_access_time` - (Optional, String) Last access time.
* `life_time` - (Optional, Int) Life cycle.
* `partition_keys` - (Optional, List) Partition keys.
* `partitions` - (Optional, List) Partitions.
* `record_count` - (Optional, Int) Record count.
* `sds` - (Optional, List) Storage object.
* `storage_size` - (Optional, Int) Storage size.
* `struct_update_time` - (Optional, String) Structure update time.
* `type` - (Optional, String) Table type: EXTERNAL_TABLE, VIRTUAL_VIEW, MATERIALIZED_VIEW.
* `view_expanded_text` - (Optional, String) View expanded text.
* `view_original_text` - (Optional, String) View original text.

The `asset` object supports the following:

* `biz_params` - (Optional, List) Additional business attributes.
* `catalog` - (Optional, String) Data catalog.
* `create_time` - (Optional, String) Create time.
* `data_version` - (Optional, Int) Data version.
* `datasource_id` - (Optional, Int) Data source primary key.
* `description` - (Optional, String) Description.
* `guid` - (Optional, String) Object GUID value.
* `id` - (Optional, Int) Primary key.
* `modified_time` - (Optional, String) Modified time.
* `name` - (Optional, String) Name.
* `owner_account` - (Optional, String) Object owner account.
* `owner` - (Optional, String) Object owner.
* `params` - (Optional, List) Additional attributes.
* `perm_values` - (Optional, List) Permissions.

The `biz_params` object of `asset` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `biz_params` object of `columns` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `biz_params` object of `dms_cols` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `biz_params` object of `partition_keys` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `columns` object supports the following:

* `biz_params` - (Optional, List) Business parameters.
* `description` - (Optional, String) Description.
* `is_partition` - (Optional, Bool) Whether partition.
* `name` - (Optional, String) Name.
* `params` - (Optional, List) Additional parameters.
* `position` - (Optional, Int) Position.
* `type` - (Optional, String) Type.

The `dms_cols` object of `sds` supports the following:

* `biz_params` - (Optional, List) Business parameters.
* `description` - (Optional, String) Description.
* `is_partition` - (Optional, Bool) Whether partition.
* `name` - (Optional, String) Name.
* `params` - (Optional, List) Additional parameters.
* `position` - (Optional, Int) Position.
* `type` - (Optional, String) Type.

The `env_props` object supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `params` object of `asset` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `params` object of `columns` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `params` object of `dms_cols` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `params` object of `partition_keys` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `params` object of `partitions` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `params` object of `sds` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `partition_keys` object supports the following:

* `biz_params` - (Optional, List) Business parameters.
* `description` - (Optional, String) Description.
* `is_partition` - (Optional, Bool) Whether partition.
* `name` - (Optional, String) Name.
* `params` - (Optional, List) Additional parameters.
* `position` - (Optional, Int) Position.
* `type` - (Optional, String) Type.

The `partitions` object supports the following:

* `create_time` - (Optional, String) Create time.
* `data_version` - (Optional, Int) Data version.
* `database_name` - (Optional, String) Database name.
* `datasource_connection_name` - (Optional, String) Data source connection name.
* `last_access_time` - (Optional, String) Last access time.
* `modified_time` - (Optional, String) Modified time.
* `name` - (Optional, String) Partition name.
* `params` - (Optional, List) Additional attributes.
* `record_count` - (Optional, Int) Record count.
* `schema_name` - (Optional, String) Schema name.
* `sds` - (Optional, List) Storage object.
* `storage_size` - (Optional, Int) Storage size.
* `table_name` - (Optional, String) Table name.
* `values` - (Optional, List) Value list.

The `perm_values` object of `asset` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `sds` object of `partitions` supports the following:

* `bucket_cols` - (Optional, List) Bucket columns.
* `compressed` - (Optional, Bool) Whether compressed.
* `dms_cols` - (Optional, List) Columns.
* `input_format` - (Optional, String) Input format.
* `location` - (Optional, String) Storage location.
* `num_buckets` - (Optional, Int) Bucket count.
* `output_format` - (Optional, String) Output format.
* `params` - (Optional, List) Additional parameters.
* `serde_lib` - (Optional, String) Serde lib.
* `serde_name` - (Optional, String) Serde name.
* `serde_params` - (Optional, List) Serde parameters.
* `sort_cols` - (Optional, List) Column sort (Expired).
* `sort_columns` - (Optional, List) Column sort fields.
* `stored_as_sub_directories` - (Optional, Bool) Whether has sub directories.

The `sds` object supports the following:

* `bucket_cols` - (Optional, List) Bucket columns.
* `compressed` - (Optional, Bool) Whether compressed.
* `dms_cols` - (Optional, List) Columns.
* `input_format` - (Optional, String) Input format.
* `location` - (Optional, String) Storage location.
* `num_buckets` - (Optional, Int) Bucket count.
* `output_format` - (Optional, String) Output format.
* `params` - (Optional, List) Additional parameters.
* `serde_lib` - (Optional, String) Serde lib.
* `serde_name` - (Optional, String) Serde name.
* `serde_params` - (Optional, List) Serde parameters.
* `sort_cols` - (Optional, List) Column sort (Expired).
* `sort_columns` - (Optional, List) Column sort fields.
* `stored_as_sub_directories` - (Optional, Bool) Whether has sub directories.

The `serde_params` object of `sds` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `sort_cols` object of `sds` supports the following:

* `col` - (Optional, String) Column name.
* `order` - (Optional, Int) Order.

The `sort_columns` object of `sds` supports the following:

* `col` - (Optional, String) Column name.
* `order` - (Optional, Int) Order.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `retention` - Hive retention version.
* `schema_name` - Schema name.


## Import

DLC DMS table can be imported using the db_name#name, e.g.

```
terraform import tencentcloud_dlc_dms_table.example tf_example_db#tf_example_table
```

