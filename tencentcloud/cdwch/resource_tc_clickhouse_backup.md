Provides a resource to open clickhouse backup

Example Usage

```hcl
resource "tencentcloud_clickhouse_backup" "backup" {
  instance_id = "cdwch-xxxxxx"
  cos_bucket_name = "xxxxxx"
}
```

Import

clickhouse backup can be imported using the id, e.g.

```
terraform import tencentcloud_clickhouse_backup.backup instance_id
```