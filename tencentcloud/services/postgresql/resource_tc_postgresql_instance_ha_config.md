Provides a resource to set postgresql instance syncMode

Example Usage

If `sync_mode` is `Semi-sync`

```hcl
resource "tencentcloud_postgresql_instance_ha_config" "example" {
  instance_id              = "postgres-gzg9jb2n"
  sync_mode                = "Semi-sync"
  max_standby_latency      = 10737418240
  max_standby_lag          = 10
  max_sync_standby_latency = 52428800
  max_sync_standby_lag     = 5
}
```

If `sync_mode` is `Async`

```hcl
resource "tencentcloud_postgresql_instance_ha_config" "example" {
  instance_id              = "postgres-gzg9jb2n"
  sync_mode                = "Async"
  max_standby_latency      = 10737418240
  max_standby_lag          = 10
}
```
