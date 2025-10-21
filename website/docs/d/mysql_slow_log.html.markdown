---
subcategory: "TencentDB for MySQL(cdb)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mysql_slow_log"
sidebar_current: "docs-tencentcloud-datasource-mysql_slow_log"
description: |-
  Use this data source to query detailed information of mysql slow_log
---

# tencentcloud_mysql_slow_log

Use this data source to query detailed information of mysql slow_log

## Example Usage

```hcl
data "tencentcloud_mysql_slow_log" "slow_log" {
  instance_id = "cdb-fitq5t9h"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance ID, in the format: cdb-c1nl9rpv. Same instance ID as displayed in the ApsaraDB for Console page.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - Details of slow query logs that meet the query conditions.
  * `date` - Backup snapshot time, time format: 2016-03-17 02:10:37.
  * `internet_url` - External network download address.
  * `intranet_url` - Intranet download address.
  * `name` - backup file name.
  * `size` - Backup file size, unit: Byte.
  * `type` - Log specific type, possible values: slowlog - slow log.


