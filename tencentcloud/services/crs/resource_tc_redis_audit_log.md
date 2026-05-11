Provides a resource to manage Redis audit log configuration.

Example Usage

```hcl
resource "tencentcloud_redis_audit_log" "example" {
  instance_id          = "crs-6eqwe3lt"
  log_sub_type         = "all"
  log_expire_day       = 7
  high_log_expire_day  = 7
  degrade_strategy     = 500
}
```

Import

Redis audit log can be imported using the instance id, e.g.

```
terraform import tencentcloud_redis_audit_log.example crs-6eqwe3lt
```
