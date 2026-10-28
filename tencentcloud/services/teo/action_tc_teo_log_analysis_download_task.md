Provides an action to create a TEO (TencentCloud EdgeOne) log analysis download task. The action creates the download task, waits for the logs to be packed, and reports the download URL of the generated file. Download tasks are retained for 3 days after creation and cannot be deleted or modified.

~> **NOTE:** Actions are supported in HashiCorp Terraform version 1.14 and later.

Example Usage

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
