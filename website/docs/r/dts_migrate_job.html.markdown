---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_migrate_job"
sidebar_current: "docs-tencentcloud-resource-dts_migrate_job"
description: |-
  Provides a resource to create a DTS migrate job
---

# tencentcloud_dts_migrate_job

Provides a resource to create a DTS migrate job

## Example Usage

```hcl
resource "tencentcloud_mysql_instance" "example" {
  instance_name     = "tf-example"
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "PassWord@123"
  slave_deploy_mode = 0
  slave_sync_mode   = 1
  availability_zone = "ap-guangzhou-7"
  mem_size          = 128000
  volume_size       = 250
  vpc_id            = "vpc-i5yyodl9"
  subnet_id         = "subnet-hhi88a58"
  intranet_port     = 3306
  security_groups   = ["sg-e6a8xxib"]
  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }
  tags = {
    createBy = "Terraform"
  }
}

resource "tencentcloud_cynosdb_cluster" "example" {
  cluster_name                 = "tf-example"
  db_mode                      = "NORMAL"
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  port                         = 3306
  password                     = "Password@123"
  force_delete                 = true
  available_zone               = "ap-guangzhou-6"
  slave_zone                   = "ap-guangzhou-7"
  vpc_id                       = "vpc-i5yyodl9"
  subnet_id                    = "subnet-hhi88a58"
  instance_cpu_core            = 2
  instance_memory_size         = 4
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 3600
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  rw_group_sg = ["sg-e6a8xxib"]
  ro_group_sg = ["sg-e6a8xxib"]
}

resource "tencentcloud_dts_migrate_service" "example" {
  src_database_type = "mysql"
  dst_database_type = "cynosdbmysql"
  src_region        = "ap-guangzhou"
  dst_region        = "ap-guangzhou"
  instance_class    = "small"
  job_name          = "tf-example"
  tags {
    tag_key   = "createBy"
    tag_value = "Terraform"
  }
}

resource "tencentcloud_dts_migrate_job" "example" {
  service_id                    = tencentcloud_dts_migrate_service.example.id
  run_mode                      = "immediate"
  auto_retry_time_range_minutes = 0
  migrate_option {
    database_table {
      object_mode = "partial"
      databases {
        db_name    = "db_name"
        db_mode    = "partial"
        table_mode = "partial"
        tables {
          table_name      = "table_name"
          new_table_name  = "new_table_name"
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
      user        = "root"
      password    = "Password@123"
      instance_id = tencentcloud_mysql_instance.example.id
    }
  }

  dst_info {
    region        = "ap-guangzhou"
    access_type   = "cdb"
    database_type = "cynosdbmysql"
    node_type     = "simple"
    info {
      user        = "root"
      password    = "Password@123"
      instance_id = tencentcloud_cynosdb_cluster.example.id
    }
  }
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

DTS migrate job can be imported using the id, e.g.

```
terraform import tencentcloud_dts_migrate_job.example dts-iy98oxba
```

