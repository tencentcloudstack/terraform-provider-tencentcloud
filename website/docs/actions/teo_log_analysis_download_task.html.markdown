---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_log_analysis_download_task"
sidebar_current: "docs-tencentcloud-action-teo_log_analysis_download_task"
description: |-
  Provides an action to create a TEO (TencentCloud EdgeOne) log analysis download task. The action creates the download task, waits for the logs to be packed, and reports the download URL of the generated file. Download tasks are retained for 3 days after creation and cannot be deleted or modified.
---

# tencentcloud_teo_log_analysis_download_task

Provides an action to create a TEO (TencentCloud EdgeOne) log analysis download task. The action creates the download task, waits for the logs to be packed, and reports the download URL of the generated file. Download tasks are retained for 3 days after creation and cannot be deleted or modified.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

## Example Usage

```hcl
action "tencentcloud_teo_log_analysis_download_task" "example" {
  config {
    zone_id    = "zone-3fkff38fyw8s"
    area       = "mainland"
    start_time = "2020-04-29T00:00:00Z"
    end_time   = "2020-04-30T00:00:00Z"
    log_type   = "l7-access-logs"
    sort       = "desc"
  }
}
```

## Argument Reference

The following arguments are supported:

* `area` - (Required, String) Data region. Valid values: `mainland`, `overseas`. For a site serving global zones, query `mainland` and `overseas` separately.
* `end_time` - (Required, String) End time of the logs to download, for example `2020-04-30T00:00:00Z`.
* `start_time` - (Required, String) Start time of the logs to download, for example `2020-04-29T00:00:00Z`. The maximum span between start time and end time is 31 days.
* `zone_id` - (Required, String) Site ID.
* `condition` - (Optional, String) Log match condition, maximum length 12KB.
* `format` - (Optional, String) File format. Valid values: `csv`. Default: `csv`.
* `log_type` - (Optional, String) Log type. Valid values: `l7-access-logs` (L7 access logs), `web-attack` (managed rules logs). Default: `l7-access-logs`.
* `sort` - (Optional, String) Time sort of the raw logs. Valid values: `asc`, `desc`. Default: `desc`.


