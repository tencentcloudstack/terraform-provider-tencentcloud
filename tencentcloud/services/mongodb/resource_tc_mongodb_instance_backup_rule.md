Provides a resource to create mongodb instance backup rule

Example Usage

```hcl
resource "tencentcloud_mongodb_instance_backup_rule" "example" {
  instance_id             = "cmgo-rnht8d3d"
  backup_method           = 0
  backup_time             = 10
  backup_retention_period = 7
  backup_version          = 1
}
```

Import

mongodb instance backup rule can be imported using the id, e.g.

```
terraform import tencentcloud_mongodb_instance_backup_rule.example cmgo-rnht8d3d
```