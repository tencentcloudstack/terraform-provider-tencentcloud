---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_lineage_attachment"
sidebar_current: "docs-tencentcloud-resource-wedata_lineage_attachment"
description: |-
  Provides a resource to create a WeData lineage attachment
---

# tencentcloud_wedata_lineage_attachment

Provides a resource to create a WeData lineage attachment

~> **NOTE:** Do not use the same relation parameters for lineage binding, as this will cause overwriting.

## Example Usage

```hcl
resource "tencentcloud_wedata_lineage_attachment" "example" {
  relations {
    source {
      resource_unique_id = "2s5veseIo2AXGOHJkKjBvQ"
      resource_type      = "TABLE"
      platform           = "WEDATA"
      resource_name      = "db_demo.1"
      description        = "DLC"
    }

    target {
      resource_unique_id = "fM8OgzE-AM2h4aaJmdXoPg"
      resource_type      = "TABLE"
      platform           = "WEDATA"
      resource_name      = "db_demo.2"
      description        = "DLC"
    }

    processes {
      process_id       = "20241107221758402"
      process_type     = "SCHEDULE_TASK"
      platform         = "WEDATA"
      process_sub_type = "SQL_TASK"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `relations` - (Required, List, ForceNew) List of lineage relationships to be registered.

The `process_properties` object of `processes` supports the following:

* `name` - (Optional, String, ForceNew) Property name.
* `value` - (Optional, String, ForceNew) Property value.

The `processes` object of `relations` supports the following:

* `platform` - (Required, String, ForceNew) WEDATA, THIRD.
* `process_id` - (Required, String, ForceNew) Original unique ID.
* `process_type` - (Required, String, ForceNew) Task type.
    //Scheduled task
    SCHEDULE_TASK,
    //Integration task
    INTEGRATION_TASK,
    //Third-party reporting
    THIRD_REPORT,
    //Data modeling
    TABLE_MODEL,
    //Model creates metric
    MODEL_METRIC,
    //Atomic metric creates derived metric
    METRIC_METRIC,
    //Data service
    DATA_SERVICE.
* `lineage_node_id` - (Optional, String, ForceNew) Lineage task unique node ID.
* `process_properties` - (Optional, List, ForceNew) Additional extension parameters.
* `process_sub_type` - (Optional, String, ForceNew) Task subtype.
 SQL_TASK,
    //Integrated real-time task lineage
    INTEGRATED_STREAM,
    //Integrated offline task lineage
    INTEGRATED_OFFLINE.

The `relations` object supports the following:

* `processes` - (Required, List, ForceNew) Lineage processing process.
* `source` - (Required, List, ForceNew) Source.
* `target` - (Required, List, ForceNew) Target.

The `resource_properties` object of `source` supports the following:

* `name` - (Optional, String, ForceNew) Property name.
* `value` - (Optional, String, ForceNew) Property value.

The `resource_properties` object of `target` supports the following:

* `name` - (Optional, String, ForceNew) Property name.
* `value` - (Optional, String, ForceNew) Property value.

The `source` object of `relations` supports the following:

* `platform` - (Required, String, ForceNew) Source: WEDATA|THIRD.
Default is wedata.
* `resource_type` - (Required, String, ForceNew) Entity type.
TABLE|METRIC|MODEL|SERVICE|COLUMN.
* `resource_unique_id` - (Required, String, ForceNew) Entity original unique ID.\n
Note: When lineage is for table columns, the unique ID should be passed as TableResourceUniqueId::FieldName.
* `create_time` - (Optional, String, ForceNew) Creation time.
* `description` - (Optional, String, ForceNew) Description: table type | metric description | model description | field description.
* `lineage_node_id` - (Optional, String, ForceNew) Lineage node unique identifier.
* `resource_name` - (Optional, String, ForceNew) Business name: database.table | metric name | model name | field name.
* `resource_properties` - (Optional, List, ForceNew) Resource additional extension parameters.
* `update_time` - (Optional, String, ForceNew) Update time.

The `target` object of `relations` supports the following:

* `platform` - (Required, String, ForceNew) Source: WEDATA|THIRD.
Default is wedata.
* `resource_type` - (Required, String, ForceNew) Entity type.
TABLE|METRIC|MODEL|SERVICE|COLUMN.
* `resource_unique_id` - (Required, String, ForceNew) Entity original unique ID.\n
Note: When lineage is for table columns, the unique ID should be passed as TableResourceUniqueId::FieldName.
* `create_time` - (Optional, String, ForceNew) Creation time.
* `description` - (Optional, String, ForceNew) Description: table type | metric description | model description | field description.
* `lineage_node_id` - (Optional, String, ForceNew) Lineage node unique identifier.
* `resource_name` - (Optional, String, ForceNew) Business name: database.table | metric name | model name | field name.
* `resource_properties` - (Optional, List, ForceNew) Resource additional extension parameters.
* `update_time` - (Optional, String, ForceNew) Update time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



