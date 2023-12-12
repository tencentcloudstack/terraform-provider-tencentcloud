Provides a resource to create a mongodb instance_backup

Example Usage

```hcl
resource "tencentcloud_mongodb_instance_backup" "instance_backup" {
  instance_id = "cmgo-9d0p6umb"
  backup_method = 0
  backup_remark = "my backup"
}
```