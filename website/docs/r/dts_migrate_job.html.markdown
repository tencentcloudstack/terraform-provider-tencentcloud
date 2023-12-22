---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_migrate_job"
sidebar_current: "docs-tencentcloud-resource-dts_migrate_job"
description: |-
  Provides a resource to create a dts migrate_job
---

# tencentcloud_dts_migrate_job

Provides a resource to create a dts migrate_job

## Example Usage

```hcl
resource "tencentcloud_cynosdb_cluster" "foo" {
  available_zone               = var.availability_zone
  vpc_id                       = local.vpc_id
  subnet_id                    = local.subnet_id
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf-cynosdb-mysql"
  password                     = "cynos@123"
  instance_maintain_duration   = 3600
  instance_maintain_start_time = 10800
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  instance_cpu_core    = 1
  instance_memory_size = 2
  param_items {
    name          = "character_set_server"
    current_value = "utf8"
  }
  param_items {
    name          = "time_zone"
    current_value = "+09:00"
  }
  param_items {
    name          = "lower_case_table_names"
    current_value = "1"
  }

  force_delete = true

  rw_group_sg = [
    local.sg_id
  ]
  ro_group_sg = [
    local.sg_id
  ]
  prarm_template_id = var.my_param_template
}

resource "tencentcloud_dts_migrate_service" "service" {
  src_database_type = "mysql"
  dst_database_type = "cynosdbmysql"
  src_region        = "ap-guangzhou"
  dst_region        = "ap-guangzhou"
  instance_class    = "small"
  job_name          = "tf_test_migration_service_1"
  tags {
    tag_key   = "aaa"
    tag_value = "bbb"
  }
}

resource "tencentcloud_dts_migrate_job" "job" {
  service_id = tencentcloud_dts_migrate_service.service.id
  run_mode   = "immediate"
  migrate_option {
    database_table {
      object_mode = "partial"
      databases {
        db_name    = "tf_ci_test"
        db_mode    = "partial"
        table_mode = "partial"
        tables {
          table_name      = "test"
          new_table_name  = "test_%s"
          table_edit_mode = "rename"
        }
      }
    }
  }
  src_info {
    region        = "ap-guangzhou"
    access_type   = "cdb"
    database_type = "mysql"
    node_type     = "simple"
    info {
      user        = "user_name"
      password    = "your_pw"
      instance_id = "cdb-fitq5t9h"
    }

  }
  dst_info {
    region        = "ap-guangzhou"
    access_type   = "cdb"
    database_type = "cynosdbmysql"
    node_type     = "simple"
    info {
      user        = "user_name"
      password    = "your_pw"
      instance_id = tencentcloud_cynosdb_cluster.foo.id
    }
  }
  auto_retry_time_range_minutes = 0
}

resource "tencentcloud_dts_migrate_job_start_operation" "start" {
  job_id = tencentcloud_dts_migrate_job.job.id
}
```

## Argument Reference

The following arguments are supported:

* `dst_info` - (Required, List) DstInfo.
* `migrate_option` - (Required, List) Migration job configuration options, used to describe how the task performs migration.
* `run_mode` - (Required, String) Run Mode. eg:immediate,timed.
* `service_id` - (Required, String) Migrate service Id from `tencentcloud_dts_migrate_service`.
* `src_info` - (Required, List) SrcInfo.
* `auto_retry_time_range_minutes` - (Optional, Int) AutoRetryTimeRangeMinutes.
* `expect_run_time` - (Optional, String) ExpectRunTime.

The `consistency` object of `migrate_option` supports the following:

* `mode` - (Optional, String) ConsistencyOption.

The `database_table` object of `migrate_option` supports the following:

* `object_mode` - (Required, String) Object mode. eg:all,partial.
* `advanced_objects` - (Optional, Set) AdvancedObjects.
* `databases` - (Optional, List) The database list.

The `databases` object of `database_table` supports the following:

* `db_mode` - (Optional, String) DB selection mode:all (for all objects under the current object), partial (partial objects), when the ObjectMode is partial, this item is required.
* `db_name` - (Optional, String) database name.
* `event_mode` - (Optional, String) EventMode.
* `events` - (Optional, Set) Events.
* `function_mode` - (Optional, String) FunctionMode.
* `functions` - (Optional, Set) Functions.
* `new_db_name` - (Optional, String) New database name.
* `new_schema_name` - (Optional, String) schema name after migration or synchronization.
* `procedure_mode` - (Optional, String) ProcedureMode.
* `procedures` - (Optional, Set) Procedures.
* `role_mode` - (Optional, String) RoleMode.
* `roles` - (Optional, List) Roles.
* `schema_mode` - (Optional, String) schema mode: all,partial.
* `schema_name` - (Optional, String) schema name.
* `table_mode` - (Optional, String) table mode: all,partial.
* `tables` - (Optional, List) tables list.
* `trigger_mode` - (Optional, String) TriggerMode.
* `triggers` - (Optional, Set) Triggers.
* `view_mode` - (Optional, String) ViewMode.
* `views` - (Optional, List) Views.

The `dst_info` object supports the following:

* `access_type` - (Required, String) AccessType.
* `database_type` - (Required, String) DatabaseType.
* `info` - (Required, List) Info.
* `node_type` - (Required, String) NodeType.
* `region` - (Required, String) Region.
* `extra_attr` - (Optional, List) ExtraAttr.
* `supplier` - (Optional, String) Supplier.

The `extra_attr` object of `dst_info` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `extra_attr` object of `migrate_option` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `extra_attr` object of `src_info` supports the following:

* `key` - (Optional, String) Key.
* `value` - (Optional, String) Value.

The `info` object of `dst_info` supports the following:

* `account_mode` - (Optional, String) Account Mode.
* `account_role` - (Optional, String) Account Role.
* `account` - (Optional, String) Account.
* `ccn_gw_id` - (Optional, String) CcnGwId.
* `cvm_instance_id` - (Optional, String) CvmInstanceId.
* `db_kernel` - (Optional, String) DbKernel.
* `engine_version` - (Optional, String) Engine Version.
* `host` - (Optional, String) Host.
* `instance_id` - (Optional, String) InstanceId.
* `password` - (Optional, String) Password.
* `port` - (Optional, Int) Port.
* `role` - (Optional, String) Role.
* `subnet_id` - (Optional, String) SubnetId.
* `tmp_secret_id` - (Optional, String) Tmp SecretId.
* `tmp_secret_key` - (Optional, String) Tmp SecretKey.
* `tmp_token` - (Optional, String) Tmp Token.
* `uniq_dcg_id` - (Optional, String) UniqDcgId.
* `uniq_vpn_gw_id` - (Optional, String) UniqVpnGwId.
* `user` - (Optional, String) User.
* `vpc_id` - (Optional, String) VpcId.

The `info` object of `src_info` supports the following:

* `account_mode` - (Optional, String) AccountMode.
* `account_role` - (Optional, String) AccountRole.
* `account` - (Optional, String) Account.
* `ccn_gw_id` - (Optional, String) CcnGwId.
* `cvm_instance_id` - (Optional, String) CvmInstanceId.
* `db_kernel` - (Optional, String) DbKernel.
* `engine_version` - (Optional, String) EngineVersion.
* `host` - (Optional, String) Host.
* `instance_id` - (Optional, String) InstanceId.
* `password` - (Optional, String) Password.
* `port` - (Optional, Int) Port.
* `role` - (Optional, String) Role.
* `subnet_id` - (Optional, String) SubnetId.
* `tmp_secret_id` - (Optional, String) TmpSecretId.
* `tmp_secret_key` - (Optional, String) TmpSecretKey.
* `tmp_token` - (Optional, String) TmpToken.
* `uniq_dcg_id` - (Optional, String) UniqDcgId.
* `uniq_vpn_gw_id` - (Optional, String) UniqVpnGwId.
* `user` - (Optional, String) User.
* `vpc_id` - (Optional, String) VpcId.

The `migrate_option` object supports the following:

* `database_table` - (Required, List) Migration object option, you need to tell the migration service which library table objects to migrate.
* `consistency` - (Optional, List) Consistency.
* `extra_attr` - (Optional, List) ExtraAttr.
* `is_dst_read_only` - (Optional, Bool) IsDstReadOnly.
* `is_migrate_account` - (Optional, Bool) IsMigrateAccount.
* `is_override_root` - (Optional, Bool) IsOverrideRoot.
* `migrate_type` - (Optional, String) MigrateType.

The `roles` object of `databases` supports the following:

* `new_role_name` - (Optional, String) NewRoleName.
* `role_name` - (Optional, String) RoleName.

The `src_info` object supports the following:

* `access_type` - (Required, String) AccessType.
* `database_type` - (Required, String) DatabaseType.
* `info` - (Required, List) Info.
* `node_type` - (Required, String) NodeType.
* `region` - (Required, String) Region.
* `extra_attr` - (Optional, List) ExtraAttr.
* `supplier` - (Optional, String) Supplier.

The `tables` object of `databases` supports the following:

* `new_table_name` - (Optional, String) new table name.
* `table_edit_mode` - (Optional, String) table edit mode.
* `table_name` - (Optional, String) table name.
* `tmp_tables` - (Optional, Set) temporary tables.

The `views` object of `databases` supports the following:

* `new_view_name` - (Optional, String) NewViewName.
* `view_name` - (Optional, String) ViewName.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Migrate job status.


## Import

dts migrate_job can be imported using the id, e.g.

```
terraform import tencentcloud_dts_migrate_job.migrate_job migrate_config_id
```

