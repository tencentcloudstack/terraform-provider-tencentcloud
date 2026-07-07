---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_backup_plan_config"
sidebar_current: "docs-tencentcloud-resource-postgresql_backup_plan_config"
description: |-
  Provides a resource to create a postgres backup_plan_config
---

# tencentcloud_postgresql_backup_plan_config

Provides a resource to create a postgres backup_plan_config

## Example Usage

```hcl
resource "tencentcloud_postgresql_backup_plan_config" "backup_plan_config" {
  db_instance_id               = local.pgsql_id
  min_backup_start_time        = "01:00:00"
  max_backup_start_time        = "02:00:00"
  base_backup_retention_period = 7
  backup_period                = ["monday", "wednesday", "friday"]
  log_backup_retention_period  = 7
  plan_name                    = "default"
  backup_method                = "physical"
}
```

## Argument Reference

The following arguments are supported:

* `db_instance_id` - (Required, String) instance id.
* `backup_method` - (Optional, String) Backup method. Valid values: `physical` (physical backup), `logical` (logical backup), `snapshot` (snapshot backup).
* `backup_period` - (Optional, Set: [`String`]) Backup cycle, which means on which days each week the instance will be backed up. The parameter value should be the lowercase names of the days of the week.
* `base_backup_retention_period` - (Optional, Int) Backup retention period in days. Value range:7-1830.
* `log_backup_retention_period` - (Optional, Int) Log backup retention period in days. Value range: 7-1830.
* `max_backup_start_time` - (Optional, String) The latest time to start a backup.
* `min_backup_start_time` - (Optional, String) The earliest time to start a backup.
* `plan_id` - (Optional, String) Backup plan ID, used to specify which backup plan to modify. If not set, the default backup plan will be modified.
* `plan_name` - (Optional, String) Custom name of the backup plan to modify.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

postgres backup_plan_config can be imported using the id, e.g.

```
terraform import tencentcloud_postgresql_backup_plan_config.backup_plan_config backup_plan_config_id
```

