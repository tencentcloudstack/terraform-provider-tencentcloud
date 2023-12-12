Use this data source to query detailed information of mysql backup_summaries

Example Usage

```hcl
data "tencentcloud_mysql_backup_summaries" "backup_summaries" {
  product = "mysql"
  order_by = "BackupVolume"
  order_direction = "ASC"
}
```