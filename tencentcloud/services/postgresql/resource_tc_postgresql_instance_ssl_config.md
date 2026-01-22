Provides a resource to create a postgres instance ssl config

~> **NOTE:** If `ssl_enabled` is `false`, Please do not set `connect_address` field.

Example Usage

Enable ssl config

```hcl
resource "tencentcloud_postgresql_instance_ssl_config" "example" {
  db_instance_id  = "postgres-5wux9sub"
  ssl_enabled     = true
  connect_address = "10.0.0.12"
}
```

Disable ssl config

```hcl
resource "tencentcloud_postgresql_instance_ssl_config" "example" {
  db_instance_id  = "postgres-5wux9sub"
  ssl_enabled     = false
}
```

Import

postgres instance ssl config can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_instance_ssl_config.example postgres-5wux9sub
```
