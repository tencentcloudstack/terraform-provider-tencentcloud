Provides a resource to create a CynosDB audit service

Example Usage

If audit_all is true

```hcl
resource "tencentcloud_cynosdb_audit_service" "example" {
  instance_id         = "cynosdbmysql-ins-f9j6sopi"
  log_expire_day      = 30
  high_log_expire_day = 7
  audit_all           = true
}
```

If audit_all is false

```hcl
resource "tencentcloud_cynosdb_audit_service" "example" {
  instance_id         = "cynosdbmysql-ins-f9j6sopi"
  log_expire_day      = 30
  high_log_expire_day = 7
  rule_template_ids   = ["cynosdb-art-riwq2vx0"]
  audit_all           = false
}
```

Import

CynosDB audit service can be imported using the id, e.g.

```
terraform import tencentcloud_cynosdb_audit_service.example cynosdbmysql-ins-f9j6sopi
```
