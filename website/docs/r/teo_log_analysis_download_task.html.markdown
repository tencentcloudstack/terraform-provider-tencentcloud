---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_log_analysis_download_task"
sidebar_current: "docs-tencentcloud-resource-teo_log_analysis_download_task"
description: |-
  Provides a resource to create a TEO (EdgeOne) log analysis download task.
---

# tencentcloud_teo_log_analysis_download_task

Provides a resource to create a TEO (EdgeOne) log analysis download task.

## Example Usage

```hcl
resource "tencentcloud_teo_log_analysis_download_task" "example" {
  zone_id    = "zone-3fkff38fyw8s"
  area       = "mainland"
  start_time = "2020-04-29T00:00:00Z"
  end_time   = "2020-04-30T00:00:00Z"
  log_type   = "l7-access-logs"
  format     = "csv"
  sort       = "desc"
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Required, String, ForceNew) Data region, valid values: `mainland`, `overseas`.
* `end_time` - (Required, String, ForceNew) End time, example: `2020-04-30T00:00:00Z`. The maximum time span from start time to end time for a single query is 31 days.
* `start_time` - (Required, String, ForceNew) Start time, example: `2020-04-29T00:00:00Z`.
* `zone_id` - (Required, String, ForceNew) Site ID.
* `condition` - (Optional, String, ForceNew) Log match condition, max length 12KB.
* `format` - (Optional, String, ForceNew) File format, valid values: `csv`. Default: `csv`.
* `log_type` - (Optional, String, ForceNew) Log type, valid values: `l7-access-logs` (L7 access logs), `web-attack` (managed rules logs). Default: `l7-access-logs`.
* `sort` - (Optional, String, ForceNew) Time sort of raw logs, valid values: `asc`, `desc`. Default: `desc`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Task creation time. Tasks are retained for 3 days after creation.
* `expire_time` - Download task expiration time. The download URL is unavailable after expiration.
* `status` - Task status, valid values: `loading`, `failed`, `completed`.
* `task_id` - Log analysis download task ID.
* `url` - Download URL, only returned when `status = completed`.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `10m`) Used when creating the resource.

## Import

TEO log analysis download task can be imported using the composite id `zoneId#area#taskId`, e.g.

```
terraform import tencentcloud_teo_log_analysis_download_task.example zone-3fkff38fyw8s#mainland#task-yyy
```

