Provides a resource to create mongodb instance backup rule

Example Usage

```hcl
resource "tencentcloud_mongodb_instance_backup" "backup_rule" {
  instance_id = "cmgo-xxxxxx"
  backup_method = 0
  backup_time = 10
}
```

Import

mongodb instance backup rule can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance_backup.backup_rule ${instanceId}
```