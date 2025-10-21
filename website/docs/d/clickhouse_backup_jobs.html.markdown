---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_backup_jobs"
sidebar_current: "docs-tencentcloud-datasource-clickhouse_backup_jobs"
description: |-
  Use this data source to query detailed information of clickhouse backup jobs
---

# tencentcloud_clickhouse_backup_jobs

Use this data source to query detailed information of clickhouse backup jobs

## Example Usage

```hcl
data "tencentcloud_clickhouse_backup_jobs" "backup_jobs" {
  instance_id = "cdwch-xxxxxx"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String) Instance id.
* `begin_time` - (Optional, String) Begin time.
* `end_time` - (Optional, String) End time.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `back_up_jobs` - Back up jobs.
  * `back_up_size` - Back up size.
  * `back_up_time` - Back up create time.
  * `back_up_type` - Back up type.
  * `expire_time` - Back up expire time.
  * `job_id` - Back up job id.
  * `job_status` - Job status.
  * `snapshot` - Back up job name.


