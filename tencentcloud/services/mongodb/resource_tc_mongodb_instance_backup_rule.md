Provides a resource to create mongodb backup rule

Example Usage

```hcl
resource "tencentcloud_mongodb_backup_rule" "backup_rule" {
  instance_id = "cmgo-xxxxxx"
  backup_method = 0
  backup_time = 10
}
```

Import

mongodb backup_rule can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_backup_rule.backup_rule ${instanceId}
```