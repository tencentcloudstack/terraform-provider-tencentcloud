---
subcategory: "TencentDB for PostgreSQL(PostgreSQL)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_postgresql_backup_plan"
sidebar_current: "docs-tencentcloud-resource-postgresql_backup_plan"
description: |-
  Provides a resource to create a PostgreSQL backup plan
---

# tencentcloud_postgresql_backup_plan

Provides a resource to create a PostgreSQL backup plan

## Example Usage

```hcl
resource "tencentcloud_postgresql_backup_plan" "example" {
  db_instance_id               = "postgres-ckwcgdf1"
  plan_name                    = "tf-example"
  backup_period_type           = "month"
  backup_period                = ["1", "2", "15"]
  min_backup_start_time        = "01:00:00"
  max_backup_start_time        = "03:00:00"
  base_backup_retention_period = 30
  log_backup_retention_period  = 7
  backup_method                = "physical"
}
```

## Argument Reference

The following arguments are supported:

* `backup_period_type` - (Required, String, ForceNew) Backup period type, currently only supports month.
* `backup_period` - (Required, Set: [`String`]) Backup dates, such as backing up on the 2nd of each month.
* `db_instance_id` - (Required, String, ForceNew) Instance ID.
* `plan_name` - (Required, String) Backup plan name.
* `backup_method` - (Optional, String) Backup method. Enumerated values: physical, logical, snapshot.
* `base_backup_retention_period` - (Optional, Int) Data backup retention period in days. Value range: [0, 30000).
* `log_backup_retention_period` - (Optional, Int) Log backup retention period in days. Value range: 7-1830.
* `max_backup_start_time` - (Optional, String) The latest time to start a backup.
* `min_backup_start_time` - (Optional, String) The earliest time to start a backup.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `plan_id` - Backup plan ID.


## Import

PostgreSQL backup plan can be imported using the db_instance_id#plan_id, e.g.

```
terraform import tencentcloud_postgresql_backup_plan.example postgres-ckwcgdf1#0ba458c5-2468-5804-ab50-ca556e88c6a3
```

