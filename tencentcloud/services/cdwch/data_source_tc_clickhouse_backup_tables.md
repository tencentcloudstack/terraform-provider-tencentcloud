Use this data source to query detailed information of clickhouse backup tables

Example Usage

```hcl
data "tencentcloud_clickhouse_backup_tables" "backup_tables" {
  instance_id = "cdwch-xxxxxx"
}
```