---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_list_process_lineage"
sidebar_current: "docs-tencentcloud-datasource-wedata_list_process_lineage"
description: |-
  Use this data source to query detailed information of WeData list process lineage
---

# tencentcloud_wedata_list_process_lineage

Use this data source to query detailed information of WeData list process lineage

## Example Usage

```hcl
data "tencentcloud_wedata_list_process_lineage" "example" {
  process_id   = "20241107221758402"
  process_type = "SCHEDULE_TASK"
  platform     = "WEDATA"
}
```

## Argument Reference

The following arguments are supported:

* `platform` - (Required, String) Source: WEDATA|THIRD, default WEDATA.
* `process_id` - (Required, String) Task unique ID.
* `process_type` - (Required, String) Task type: SCHEDULE_TASK, INTEGRATION_TASK, THIRD_REPORT, TABLE_MODEL, MODEL_METRIC, METRIC_METRIC, DATA_SERVICE.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Lineage pair list.
  * `processes` - Lineage processing procedures.
    * `lineage_node_id` - Lineage task unique node ID.
    * `platform` - WEDATA, THIRD.
    * `process_id` - Original unique ID.
    * `process_properties` - Additional extension parameters.
      * `name` - Property name.
      * `value` - Property value.
    * `process_sub_type` - Task subtype.
SQL_TASK,
INTEGRATED_STREAM,
INTEGRATED_OFFLINE.
    * `process_type` - Task type.
SCHEDULE_TASK,
INTEGRATION_TASK,
THIRD_REPORT,
TABLE_MODEL,
MODEL_METRIC,
METRIC_METRIC,
DATA_SERVICE.
  * `source` - Source.
    * `create_time` - Creation time.
    * `description` - Description: table type|metric description|model description|field description.
    * `lineage_node_id` - Lineage node unique identifier.
    * `platform` - Source: WEDATA|THIRD.
Default wedata.
    * `resource_name` - Business name: database.table|metric name|model name|field name.
    * `resource_properties` - Resource additional extension parameters.
      * `name` - Property name.
      * `value` - Property value.
    * `resource_type` - Entity type.
TABLE|METRIC|MODEL|SERVICE|COLUMN.
    * `resource_unique_id` - Entity original unique ID.

Note: When lineage is for table columns, the unique ID should be TableResourceUniqueId::FieldName.
    * `update_time` - Update time.
  * `target` - Target.
    * `create_time` - Creation time.
    * `description` - Description: table type|metric description|model description|field description.
    * `lineage_node_id` - Lineage node unique identifier.
    * `platform` - Source: WEDATA|THIRD.
Default wedata.
    * `resource_name` - Business name: database.table|metric name|model name|field name.
    * `resource_properties` - Resource additional extension parameters.
      * `name` - Property name.
      * `value` - Property value.
    * `resource_type` - Entity type.
TABLE|METRIC|MODEL|SERVICE|COLUMN.
    * `resource_unique_id` - Entity original unique ID.

Note: When lineage is for table columns, the unique ID should be TableResourceUniqueId::FieldName.
    * `update_time` - Update time.


