Provides a resource to recover a clickhouse back up

Example Usage

```hcl
resource "tencentcloud_clickhouse_recover_backup_job" "recover_backup_job" {
  instance_id = "cdwch-xxxxxx"
  back_up_job_id = 1234
}
```