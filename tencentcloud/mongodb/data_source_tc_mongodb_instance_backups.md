Use this data source to query detailed information of mongodb instance_backups

Example Usage

```hcl
data "tencentcloud_mongodb_instance_backups" "instance_backups" {
  instance_id = "cmgo-9d0p6umb"
  backup_method = 0
}
```