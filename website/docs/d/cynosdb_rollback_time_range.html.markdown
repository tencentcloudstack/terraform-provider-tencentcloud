---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_rollback_time_range"
sidebar_current: "docs-tencentcloud-datasource-cynosdb_rollback_time_range"
description: |-
  Use this data source to query detailed information of cynosdb rollback_time_range
---

# tencentcloud_cynosdb_rollback_time_range

Use this data source to query detailed information of cynosdb rollback_time_range

## Example Usage

```hcl
data "tencentcloud_cynosdb_rollback_time_range" "rollback_time_range" {
  cluster_id = "cynosdbmysql-bws8h88b"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String) Cluster ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `rollback_time_ranges` - Reversible time range.
  * `time_range_end` - End time.
  * `time_range_start` - start time.
* `time_range_end` - Effective regression time range end time point (obsolete) Note: This field may return null, indicating that a valid value cannot be obtained.
* `time_range_start` - Effective regression time range start time point (obsolete) Note: This field may return null, indicating that a valid value cannot be obtained.


