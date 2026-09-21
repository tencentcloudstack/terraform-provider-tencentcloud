Provides a resource to create a TEO (EdgeOne) log analysis download task.

Example Usage

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

Import

TEO log analysis download task can be imported using the composite id `zoneId#area#taskId`, e.g.

```
terraform import tencentcloud_teo_log_analysis_download_task.example zone-3fkff38fyw8s#mainland#task-yyy
```
