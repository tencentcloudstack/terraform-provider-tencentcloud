Use this data source to query detailed information of clickhouse backup job detail

Example Usage

```hcl
data "tencentcloud_clickhouse_backup_job_detail" "backup_job_detail" {
  instance_id = "cdwch-xxxxxx"
  back_up_job_id = 1234
}
```