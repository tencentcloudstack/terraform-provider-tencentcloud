---
subcategory: "Cloud Log Service(CLS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cls_dlc_deliver"
sidebar_current: "docs-tencentcloud-resource-cls_dlc_deliver"
description: |-
  Provides a resource to create CLS dcl deliver
---

# tencentcloud_cls_dlc_deliver

Provides a resource to create CLS dcl deliver

## Example Usage

```hcl
resource "tencentcloud_cls_dlc_deliver" "example" {
  topic_id     = "5ba3b3eb-7459-4807-82d9-c98236d2e100"
  name         = "tf-example"
  deliver_type = 0
  start_time   = 1775118742

  dlc_info {
    table_info {
      data_directory = "DataLakeCategary"
      database_name  = "tf_example_db"
      table_name     = "tf_example_table"
    }

    field_infos {
      cls_field      = "info"
      dlc_field      = "info"
      dlc_field_type = "string"
      disable        = false
    }

    field_infos {
      cls_field      = "int_key"
      dlc_field      = "int_key"
      dlc_field_type = "int"
      disable        = false
    }

    field_infos {
      cls_field      = "bool_key"
      dlc_field      = "bool_key"
      dlc_field_type = "boolean"
      disable        = false
    }

    field_infos {
      cls_field      = "float_key"
      dlc_field      = "float_key"
      dlc_field_type = "float"
      disable        = false
    }

    field_infos {
      cls_field      = "double_key"
      dlc_field      = "double_key"
      dlc_field_type = "double"
      disable        = false
    }


    partition_infos {
      cls_field      = "__TIMESTAMP__"
      dlc_field      = "date_key"
      dlc_field_type = "date"
    }

    partition_extra {
      time_format = "/%Y/%m/%d/%H"
      time_zone   = "UTC+08:00"
    }
  }

  max_size         = 128
  interval         = 300
  has_services_log = 2
}
```

## Argument Reference

The following arguments are supported:

* `deliver_type` - (Required, Int) Delivery type. `0`: batch delivery, `1`: real-time delivery.
* `dlc_info` - (Required, List) DLC configuration information.
* `name` - (Required, String) Task name. Length does not exceed 64 characters, starts with a letter, accepts 0-9, a-z, A-Z, _, -, Chinese characters.
* `start_time` - (Required, Int) Start time of the delivery time range (Unix timestamp).
* `topic_id` - (Required, String, ForceNew) Log topic ID.
* `end_time` - (Optional, Int) End time of the delivery time range (Unix timestamp). If empty, no time limit. Must be greater than `start_time` when set.
* `has_services_log` - (Optional, Int) Whether to enable delivery service logs. `1`: disabled, `2`: enabled. Default is enabled.
* `interval` - (Optional, Int) Delivery interval in seconds. Required when `deliver_type=0`. Range: 300 <= Interval <= 900.
* `max_size` - (Optional, Int) Delivery file size in MB. Required when `deliver_type=0`. Range: 5 <= MaxSize <= 256.
* `status` - (Optional, Int) Task status. `1`: running, `2`: stopped.

The `dlc_info` object supports the following:

* `table_info` - (Required, List) DLC table information.
* `field_infos` - (Optional, List) DLC data field information.
* `partition_extra` - (Optional, List) DLC partition extra information.
* `partition_infos` - (Optional, List) DLC partition information.

The `field_infos` object of `dlc_info` supports the following:

* `cls_field` - (Required, String) Field name in CLS log.
* `dlc_field_type` - (Required, String) DLC field type, e.g. `string`, `int`, `struct`.
* `dlc_field` - (Required, String) Column name in DLC table.
* `disable` - (Optional, Bool) Whether to disable this field.
* `fill_field` - (Optional, String) Fill field when parsing fails.

The `partition_extra` object of `dlc_info` supports the following:

* `time_format` - (Optional, String) Time format, e.g. `/%Y/%m/%d/%H`.
* `time_zone` - (Optional, String) Time zone, e.g. `UTC+08:00`.

The `partition_infos` object of `dlc_info` supports the following:

* `cls_field` - (Required, String) Field name in CLS log.
* `dlc_field_type` - (Required, String) DLC field type.
* `dlc_field` - (Required, String) Column name in DLC table.

The `table_info` object of `dlc_info` supports the following:

* `data_directory` - (Required, String) Data directory.
* `database_name` - (Required, String) Database name.
* `table_name` - (Required, String) Table name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `task_id` - Delivery task ID.


## Import

CLS dcl deliver can be imported using the id (topicId#taskId), e.g.

```
terraform import tencentcloud_cls_dlc_deliver.example 715094e3-01b0-4aeb-91f5-ee9f46a4a13c#988259ca-598f-428c-8475-cf438d05468c
```

