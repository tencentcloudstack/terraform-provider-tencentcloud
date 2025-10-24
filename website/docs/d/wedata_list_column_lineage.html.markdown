---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_list_column_lineage"
sidebar_current: "docs-tencentcloud-datasource-wedata_list_column_lineage"
description: |-
  Use this data source to query detailed information of WeData list column lineage
---

# tencentcloud_wedata_list_column_lineage

Use this data source to query detailed information of WeData list column lineage

## Example Usage

```hcl
data "tencentcloud_wedata_list_column_lineage" "example" {
  table_unique_id = "B_CRyO4-3rMvNFPH_7aTaw"
  direction       = "INPUT"
  column_name     = "example_column"
  platform        = "WEDATA"
}
```

## Argument Reference

The following arguments are supported:

* `column_name` - (Required, String) Column name.
* `direction` - (Required, String) Lineage direction INPUT|OUTPUT.
* `platform` - (Required, String) Source: WEDATA|THIRD, default WEDATA.
* `table_unique_id` - (Required, String) Table unique ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Lineage record list.
  * `relation` - Relation.
    * `processes` - Lineage processing process.
      * `lineage_node_id` - Lineage task unique node ID.
      * `platform` - WEDATA, THIRD.
      * `process_id` - Original unique ID.
      * `process_properties` - Additional extension parameters.
        * `name` - Property name.
        * `value` - Property value.
      * `process_sub_type` - Task subtype
 SQL_TASK,
    //Integration real-time task lineage
    INTEGRATED_STREAM,
    //Integration offline task lineage
    INTEGRATED_OFFLINE.
      * `process_type` - Task type
    //Scheduling task
    SCHEDULE_TASK,
    //Integration task
    INTEGRATION_TASK,
    //Third-party reporting
    THIRD_REPORT,
    //Data modeling
    TABLE_MODEL,
    //Model creates metrics
    MODEL_METRIC,
    //Atomic metric creates derived metric
    METRIC_METRIC,
    //Data service
    DATA_SERVICE.
    * `relation_id` - Relation ID.
    * `source_unique_id` - Source unique lineage ID.
    * `target_unique_id` - Target unique lineage ID.
  * `resource` - Current resource.
    * `create_time` - Creation time.
    * `description` - Description: table type|metric description|model description|field description.
    * `lineage_node_id` - Lineage node unique identifier.
    * `platform` - Source: WEDATA|THIRD
default wedata.
    * `resource_name` - Business name: database.table|metric name|model name|field name.
    * `resource_properties` - Resource additional extension parameters.
      * `name` - Property name.
      * `value` - Property value.
    * `resource_type` - Entity type
TABLE|METRIC|MODEL|SERVICE|COLUMN.
    * `resource_unique_id` - Entity original unique ID.
    * `update_time` - Update time.


