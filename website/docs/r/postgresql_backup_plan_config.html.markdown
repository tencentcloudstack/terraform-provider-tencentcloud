---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_backup_plan_config"
sidebar_current: "docs-tencentcloud-resource-postgresql_backup_plan_config"
description: |-
  Provides a resource to create a postgres backup plan config
---

# tencentcloud_postgresql_backup_plan_config

Provides a resource to create a postgres backup plan config

## Example Usage

```hcl
resource "tencentcloud_postgresql_backup_plan_config" "example" {
  db_instance_id               = "postgres-ckwcgdf1"
  min_backup_start_time        = "01:00:00"
  max_backup_start_time        = "03:00:00"
  backup_period                = ["monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"]
  base_backup_retention_period = 7
  log_backup_retention_period  = 7
  backup_method                = "physical"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String, ForceNew) instance id.
* `backup_method` - (Optional, String) Backup method. Valid values: `physical` (physical backup), `logical` (logical backup), `snapshot` (snapshot backup).
* `backup_period` - (Optional, Set: [`String`]) Backup cycle, which means on which days each week the instance will be backed up. The parameter value should be the lowercase names of the days of the week.
* `base_backup_retention_period` - (Optional, Int) Backup retention period in days. Value range:7-1830.
* `log_backup_retention_period` - (Optional, Int) Log backup retention period in days. Value range: 7-1830.
* `max_backup_start_time` - (Optional, String) The latest time to start a backup.
* `min_backup_start_time` - (Optional, String) The earliest time to start a backup.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

postgres backup plan config can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_backup_plan_config.example postgres-ckwcgdf1
```

