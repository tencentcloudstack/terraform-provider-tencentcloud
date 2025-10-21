---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_recovery_time"
sidebar_current: "docs-tencentcloud-datasource-postgresql_recovery_time"
description: |-
  Use this data source to query detailed information of postgresql recovery_time
---

# tencentcloud_postgresql_recovery_time

Use this data source to query detailed information of postgresql recovery_time

## Example Usage

```hcl
data "tencentcloud_postgresql_recovery_time" "recovery_time" {
  db_instance_id = local.pgsql_id
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String) Instance ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `recovery_begin_time` - The earliest restoration time (UTC+8).
* `recovery_end_time` - The latest restoration time (UTC+8).


