---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_describe_instance_slow_queries"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_describe_instance_slow_queries"
description: |-
  Use this data source to query detailed information of cynosdb describe_instance_slow_queries
---

# tencentcloud_cynosdb_describe_instance_slow_queries

Use this data source to query detailed information of cynosdb describe_instance_slow_queries

## Example Usage

```hcl
data "tencentcloud_cynosdb_describe_instance_slow_queries" "describe_instance_slow_queries" {
  cluster_id = "cynosdbmysql-bws8h88b"
  start_time = "2023-06-01 12:00:00"
  end_time   = "2023-06-19 14:00:00"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `end_time` - (Optional, String) End time.
* `result_output_file` - (Optional, String) Used to save results.
* `start_time` - (Optional, String) start time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `binlogs` - Note to the Binlog list: This field may return null, indicating that a valid value cannot be obtained.
  * `binlog_id` - Binlog file ID.
  * `file_name` - Binlog file name.
  * `file_size` - File size in bytes.
  * `finish_time` - Latest transaction time.
  * `start_time` - Earliest transaction time.


