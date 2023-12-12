Provides a resource to delete a clickhouse back up data

Example Usage

```hcl
resource "tencentcloud_clickhouse_delete_backup_data" "delete_back_up_data" {
  instance_id = "cdwch-xxxxxx"
  back_up_job_id = 1234
}
```