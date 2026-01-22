---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_updatable_data_engines"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_updatable_data_engines"
description: |-
  Use this data source to query detailed information of DLC describe updatable data engines
---

# tencentcloud_dlc_describe_updatable_data_engines

Use this data source to query detailed information of DLC describe updatable data engines

## Example Usage

```hcl
data "tencentcloud_dlc_describe_updatable_data_engines" "example" {
  data_engine_config_command = "UpdateSparkSQLLakefsPath"
}
```

### Or

```hcl
data "tencentcloud_dlc_describe_updatable_data_engines" "example" {
  data_engine_config_command = "UpdateSparkSQLResultPath"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_config_command` - (Required, String) Operation commands of engine configuration. UpdateSparkSQLLakefsPath updates the path of managed tables, and UpdateSparkSQLResultPath updates the path of result buckets.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data_engine_basic_infos` - Basic cluster information.
  * `app_id` - User ID.
  * `create_time` - Create time.
  * `data_engine_id` - Engine ID.
  * `data_engine_name` - DataEngine name.
  * `data_engine_type` - Engine types, and the valid values are PrestoSQL, SparkSQL, and SparkBatch.
  * `message` - Returned information.
  * `state` - EData engine status: -2: deleted; -1: failed; 0: initializing; 1: suspended; 2: running; 3: ready to delete; 4: deleting.
  * `update_time` - Update time.
  * `user_uin` - Account uin.


