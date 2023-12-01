---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_updatable_data_engines"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_updatable_data_engines"
description: |-
  Use this data source to query detailed information of dlc describe_updatable_data_engines
---

# tencentcloud_dlc_describe_updatable_data_engines

Use this data source to query detailed information of dlc describe_updatable_data_engines

## Example Usage

```hcl
data "tencentcloud_dlc_describe_updatable_data_engines" "describe_updatable_data_engines" {
  data_engine_config_command = "UpdateSparkSQLLakefsPath"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_config_command` - (Required, String) Engine configuration operation command, UpdateSparkSQLLakefsPath updates the managed table path, UpdateSparkSQLResultPath updates the result bucket path.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data_engine_basic_infos` - Engine basic information.
  * `app_id` - User unique ID.
  * `create_time` - Create time.
  * `data_engine_id` - Engine unique id.
  * `data_engine_name` - Engine name.
  * `data_engine_type` - Engine type, valid values: PrestoSQL/SparkSQL/SparkBatch.
  * `message` - Returned messages.
  * `state` - Engine state, only support: 0:Init/-1:Failed/-2:Deleted/1:Pause/2:Running/3:ToBeDelete/4:Deleting.
  * `update_time` - Update time.
  * `user_uin` - User unique uin.


