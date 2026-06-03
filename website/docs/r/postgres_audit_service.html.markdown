---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgres_audit_service"
sidebar_current: "docs-tencentcloud-resource-postgres_audit_service"
description: |-
  Provides a resource to manage PostgreSQL audit service
---

# tencentcloud_postgres_audit_service

Provides a resource to manage PostgreSQL audit service

## Example Usage

```hcl
resource "tencentcloud_postgres_audit_service" "example" {
  instance_id        = "postgres-ckwcgdf1"
  log_expire_day     = 30
  hot_log_expire_day = 7
  audit_type         = "complex"
}
```

## Argument Reference

The following arguments are supported:

* `audit_type` - (Required, String) Audit type. Valid values: complex (fine-grained audit), simple (fast audit).
* `hot_log_expire_day` - (Required, Int) Hot log retention days. Valid values: 7, 30, 90, 180, 365, 1095, 1825.
* `instance_id` - (Required, String, ForceNew) PostgreSQL instance ID.
* `log_expire_day` - (Required, Int) Log retention days. Valid values: 7, 30, 90, 180, 365, 1095, 1825.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `audit_status` - Audit status. Values: ON, OFF.
* `cold_log_expire_day` - Cold log retention days.
* `cold_log_size` - Cold log size in MB.
* `hot_log_size` - Hot log size in MB.


## Import

PostgreSQL audit service can be imported using the instance_id, e.g.

```
terraform import tencentcloud_postgres_audit_service.example postgres-ckwcgdf1
```

