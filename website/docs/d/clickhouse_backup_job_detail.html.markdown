---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_backup_job_detail"
sidebar_current: "docs-tencentcloud-datasource-clickhouse_backup_job_detail"
description: |-
  Use this data source to query detailed information of clickhouse backup job detail
---

# tencentcloud_clickhouse_backup_job_detail

Use this data source to query detailed information of clickhouse backup job detail

## Example Usage

```hcl
data "tencentcloud_clickhouse_backup_job_detail" "backup_job_detail" {
  instance_id    = "cdwch-xxxxxx"
  back_up_job_id = 1234
}
```

## Argument Reference

The following arguments are supported:

* `back_up_job_id` - (Required, Int) Back up job id.
* `instance_id` - (Required, String) Instance id.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `table_contents` - Back up tables.
  * `database` - Database.
  * `ips` - Ips.
  * `rip` - Ip address of cvm.
  * `table` - Table.
  * `total_bytes` - Total bytes.
  * `v_cluster` - Virtual cluster.
  * `zoo_path` - ZK path.


