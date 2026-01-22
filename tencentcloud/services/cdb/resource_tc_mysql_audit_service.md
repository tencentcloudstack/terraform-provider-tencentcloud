Provides a resource to create a Mysql audit service

Example Usage

If audit_all is true

```hcl
resource "tencentcloud_mysql_audit_service" "example" {
  instance_id         = "cdb-3kwa3gfj"
  log_expire_day      = 30
  high_log_expire_day = 7
  audit_all           = true
}
```

If audit_all is false

```hcl
resource "tencentcloud_mysql_audit_service" "example" {
  instance_id         = "cdb-3kwa3gfj"
  log_expire_day      = 30
  high_log_expire_day = 7
  rule_template_ids   = [
    "cdb-art-3a9ww0oj"
  ]
  audit_all           = false
}
```

Import

Mysql audit service can be imported using the id, e.g.

```
terraform import tencentcloud_mysql_audit_service.example cdb-3kwa3gfj
```
