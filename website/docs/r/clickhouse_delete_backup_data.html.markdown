---
subcategory: "ClickHouse(CDWCH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_clickhouse_delete_backup_data"
sidebar_current: "docs-tencentcloud-resource-clickhouse_delete_backup_data"
description: |-
  Provides a resource to delete a clickhouse back up data
---

# tencentcloud_clickhouse_delete_backup_data

Provides a resource to delete a clickhouse back up data

## Example Usage

```hcl
resource "tencentcloud_clickhouse_delete_backup_data" "delete_back_up_data" {
  instance_id    = "cdwch-xxxxxx"
  back_up_job_id = 1234
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Instance id.
* `back_up_job_id` - (Optional, Int, ForceNew) Back up job id.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



