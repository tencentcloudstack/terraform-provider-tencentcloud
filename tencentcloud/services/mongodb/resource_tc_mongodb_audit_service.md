Provides a resource to manage MongoDB audit service configuration.

Example Usage

Full audit mode

```hcl
resource "tencentcloud_mongodb_audit_service" "example" {
  instance_id    = "cmgo-xxxxxx"
  log_expire_day = 30
  audit_all      = true
}
```

Rule-based audit mode

```hcl
resource "tencentcloud_mongodb_audit_service" "example" {
  instance_id    = "cmgo-xxxxxx"
  log_expire_day = 30
  audit_all      = false

  rule_filters {
    type    = "DB"
    compare = "EQ"
    value   = ["testdb"]
  }

  rule_filters {
    type    = "User"
    compare = "EQ"
    value   = ["admin"]
  }
}
```

Import

MongoDB audit service can be imported using the instance id, e.g.

```
terraform import tencentcloud_mongodb_audit_service.example cmgo-xxxxxx
```
