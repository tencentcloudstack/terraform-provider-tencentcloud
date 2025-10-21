---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_rollback_range_time"
sidebar_current: "docs-tencentcloud-datasource-mysql_rollback_range_time"
description: |-
  Use this data source to query detailed information of mysql rollback_range_time
---

# tencentcloud_mysql_rollback_range_time

Use this data source to query detailed information of mysql rollback_range_time

## Example Usage

```hcl
data "tencentcloud_mysql_rollback_range_time" "rollback_range_time" {
  instance_ids = ["cdb-fitq5t9h"]
}
```

## Argument Reference

The following arguments are supported:

* `instance_ids` - (Required, Set: [`String`]) A list of instance IDs, the format of a single instance ID is: cdb-c1nl9rpv. Same instance ID as displayed in the ApsaraDB for Console page.
* `backup_region` - (Optional, String) If the clone instance is not in the same region as the source instance, fill in the region where the clone instance is located, for example: ap-guangzhou.
* `is_remote_zone` - (Optional, String) Whether the clone instance is in the same zone as the source instance, yes: `false`, no: `true`.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Returned parameter information.
  * `code` - Query database error code.
  * `instance_id` - A list of instance IDs. The format of a single instance ID is: cdb-c1nl9rpv. Same as the instance ID displayed in the cloud database console page.
  * `message` - Query database error information.
  * `times` - Retrievable time range.
    * `begin` - Instance rollback start time, time format: 2016-10-29 01:06:04.
    * `end` - End time of instance rollback, time format: 2016-11-02 11:44:47.


