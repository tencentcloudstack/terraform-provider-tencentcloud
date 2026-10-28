Provides a resource to manage PostgreSQL audit service

Example Usage

```hcl
resource "tencentcloud_postgres_audit_service" "example" {
  instance_id        = "postgres-ckwcgdf1"
  log_expire_day     = 30
  hot_log_expire_day = 7
  audit_type         = "complex"
}
```

Import

PostgreSQL audit service can be imported using the instance_id, e.g.

```
terraform import tencentcloud_postgres_audit_service.example postgres-ckwcgdf1
```
