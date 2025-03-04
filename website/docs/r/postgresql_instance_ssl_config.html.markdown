---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_instance_ssl_config"
sidebar_current: "docs-tencentcloud-resource-postgresql_instance_ssl_config"
description: |-
  Provides a resource to create a postgres instance ssl config
---

# tencentcloud_postgresql_instance_ssl_config

Provides a resource to create a postgres instance ssl config

~> **NOTE:** If `ssl_enabled` is `false`, Please do not set `connect_address` field.

## Example Usage

### Enable ssl config

```hcl
resource "tencentcloud_postgresql_instance_ssl_config" "example" {
  db_instance_id  = "postgres-5wux9sub"
  ssl_enabled     = true
  connect_address = "10.0.0.12"
}
```

### Disable ssl config

```hcl
resource "tencentcloud_postgresql_instance_ssl_config" "example" {
  db_instance_id = "postgres-5wux9sub"
  ssl_enabled    = false
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) Postgres instance ID.
* `ssl_enabled` - (Required, Bool) Enable or disable SSL. true: enable; false: disable.
* `connect_address` - (Optional, String) The unique connection address protected by SSL certificate, which can be set as the internal and external IP address if it is the primary instance; If it is a read-only instance, it can be set as the instance IP or read-only group IP. This parameter is mandatory when enabling SSL or modifying SSL protected connection addresses; When SSL is turned off, this parameter will be ignored.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `ca_url` - Cloud root certificate download link.


## Import

postgres instance ssl config can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_instance_ssl_config.example postgres-5wux9sub
```

