---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_backup_tables"
sidebar_current: "docs-tencentcloud-datasource-clickhouse_backup_tables"
description: |-
  Use this data source to query detailed information of clickhouse backup tables
---

# tencentcloud_clickhouse_backup_tables

Use this data source to query detailed information of clickhouse backup tables

## Example Usage

```hcl
data "tencentcloud_clickhouse_backup_tables" "backup_tables" {
  instance_id = "cdwch-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `available_tables` - Available tables.
  * `database` - Database.
  * `ips` - Table ips.
  * `rip` - Ip address of cvm.
  * `table` - Table.
  * `total_bytes` - Table total bytes.
  * `v_cluster` - Virtual cluster.
  * `zoo_path` - Zk path.


