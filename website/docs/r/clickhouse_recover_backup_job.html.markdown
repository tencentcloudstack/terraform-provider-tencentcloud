---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_recover_backup_job"
sidebar_current: "docs-tencentcloud-resource-clickhouse_recover_backup_job"
description: |-
  Provides a resource to recover a clickhouse back up
---

# tencentcloud_clickhouse_recover_backup_job

Provides a resource to recover a clickhouse back up

## Example Usage

```hcl
resource "tencentcloud_clickhouse_recover_backup_job" "recover_backup_job" {
  instance_id    = "cdwch-xxxxxx"
  back_up_job_id = 1234
}
```

## Argument Reference

The following arguments are supported:

* `back_up_job_id` - (Required, Int, ForceNew) Back up job id.
* `instance_id` - (Required, String, ForceNew) Instance id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



