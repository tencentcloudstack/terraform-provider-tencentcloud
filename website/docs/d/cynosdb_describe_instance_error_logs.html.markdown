---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_describe_instance_error_logs"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_describe_instance_error_logs"
description: |-
  Use this data source to query detailed information of cynosdb describe_instance_error_logs
---

# tencentcloud_cynosdb_describe_instance_error_logs

Use this data source to query detailed information of cynosdb describe_instance_error_logs

## Example Usage

```hcl
data "tencentcloud_cynosdb_describe_instance_error_logs" "describe_instance_error_logs" {
  instance_id   = "cynosdbmysql-ins-afqx1hy0"
  start_time    = "2023-06-01 15:04:05"
  end_time      = "2023-06-19 15:04:05"
  order_by      = "Timestamp"
  order_by_type = "DESC"
  log_levels    = ["note", "warning"]
  key_words     = ["Aborted"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance Id.
* `end_time` - (Optional, String) End time.
* `key_words` - (Optional, Set: [`String`]) Keywords, supports fuzzy search.
* `log_levels` - (Optional, Set: [`String`]) Log levels, including error, warning, and note, support simultaneous search of multiple levels.
* `order_by_type` - (Optional, String) Sort type, with ASC and DESC enumeration values.
* `order_by` - (Optional, String) Sort fields with Timestamp enumeration values.
* `result_output_file` - (Optional, String) Used to save results.
* `start_time` - (Optional, String) start time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `error_logs` - Error log list note: This field may return null, indicating that a valid value cannot be obtained.
  * `content` - Note to log content: This field may return null, indicating that a valid value cannot be obtained.
  * `level` - Log level note: This field may return null, indicating that a valid value cannot be obtained.
  * `timestamp` - Log timestamp note: This field may return null, indicating that a valid value cannot be obtained.


