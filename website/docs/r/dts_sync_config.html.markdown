---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_sync_config"
sidebar_current: "docs-tencentcloud-resource-dts_sync_config"
description: |-
  Provides a resource to create a DTS sync config
---

# tencentcloud_dts_sync_config

Provides a resource to create a DTS sync config

## Example Usage

### Sync MySQL database to CynosDB through cdb access type

```hcl
resource "tencentcloud_mysql_instance" "example" {
  instance_name     = "tf-example"
  internet_service  = 1
  engine_version    = "5.7"
  charge_type       = "POSTPAID"
  root_password     = "Mysql@2026"
  slave_deploy_mode = 0
  slave_sync_mode   = 0
  device_type       = "CLOUD_NATIVE_CLUSTER"
  availability_zone = "ap-guangzhou-6"
  cpu               = 2
  mem_size          = 4000
  volume_size       = 200
  vpc_id            = "vpc-i5yyodl9"
  subnet_id         = "subnet-hhi88a58"
  intranet_port     = 3306
  security_groups   = ["sg-4rd5741x"]
  parameters = {
    character_set_server = "utf8"
    max_connections      = "1000"
  }

  tags = {
    createBy = "terraform"
  }

  cluster_topology {
    read_write_node {
      zone = "ap-guangzhou-6"
    }
  }

  timeouts {
    create = "30m"
    delete = "30m"
  }
}

resource "tencentcloud_cynosdb_cluster" "example" {
  available_zone               = "ap-guangzhou-6"
  vpc_id                       = "vpc-i5yyodl9"
  subnet_id                    = "subnet-hhi88a58"
  db_mode                      = "NORMAL"
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  port                         = 3306
  cluster_name                 = "tf-example"
  password                     = "CynosDB@2026"
  instance_maintain_duration   = 7200
  instance_maintain_start_time = 10800
  instance_cpu_core            = 2
  instance_memory_size         = 4
  force_delete                 = true
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]

  param_items {
    name          = "character_set_server"
    current_value = "utf8mb4"
  }

  param_items {
    name          = "lower_case_table_names"
    current_value = "0"
  }

  tags = {
    createBy = "terraform"
  }
}

resource "tencentcloud_dts_sync_job" "example" {
  pay_mode          = "PostPay"
  src_database_type = "mysql"
  src_region        = "ap-guangzhou"
  dst_database_type = "cynosdbmysql"
  dst_region        = "ap-guangzhou"
  auto_renew        = 0
  instance_class    = "micro"
  tags {
    tag_key   = "key"
    tag_value = "value"
  }
}

resource "tencentcloud_dts_sync_config" "example" {
  job_id          = tencentcloud_dts_sync_job.example.job_id
  src_access_type = "cdb"
  dst_access_type = "cdb"
  job_name        = "tf_example"
  job_mode        = "liteMode"
  run_mode        = "Immediate"

  objects {
    mode = "Partial"
    databases {
      db_name        = "testDB"
      db_mode        = "Partial"
      table_mode     = "Partial"
      view_mode      = "Partial"
      procedure_mode = "Partial"
      function_mode  = "Partial"
      tables {
        table_name  = "testTable"
        column_mode = "Partial"
        columns {
          column_name = "id"
        }

        tmp_tables = [
          "_testTable_new",
          "_testTable_old",
          "_testTable_ghc",
          "_testTable_gho",
          "_testTable_del",
        ]

        table_edit_mode = "pt"
      }
    }

    advanced_objects = [
      "procedure",
      "function",
    ]
  }

  src_info {
    region      = "ap-guangzhou"
    instance_id = tencentcloud_mysql_instance.example.id
    user        = "root"
    password    = tencentcloud_mysql_instance.example.root_password
    db_name     = "testDB"
    vpc_id      = "vpc-i5yyodl9"
    subnet_id   = "subnet-hhi88a58"
  }

  dst_info {
    region      = "ap-guangzhou"
    instance_id = tencentcloud_cynosdb_cluster.example.id
    user        = "root"
    password    = tencentcloud_cynosdb_cluster.example.password
    db_name     = "testDB"
    vpc_id      = "vpc-i5yyodl9"
    subnet_id   = "subnet-hhi88a58"
  }

  auto_retry_time_range_minutes = 5
}
```

## Argument Reference

The following arguments are supported:

* `dst_access_type` - (Required, String) Target end access type, cdb (cloud database), cvm (cloud host self-built), vpc (private network), extranet (external network), vpncloud (vpn access), dcg (dedicated line access), ccn (cloud networking ), intranet (self-developed cloud), noProxy, note that the specific optional value depends on the current link.
* `job_id` - (Required, String) Synchronization instance id (i.e. identifies a synchronization job).
* `objects` - (Required, List) Synchronize database table object information.
* `src_access_type` - (Required, String) Source access type, cdb (cloud database), cvm (cloud host self-built), vpc (private network), extranet (external network), vpncloud (vpn access), dcg (dedicated line access), ccn (cloud networking ), intranet (self-developed cloud), noProxy, note that the specific optional value depends on the current link.
* `auto_retry_time_range_minutes` - (Optional, Int) The time period of automatic retry, can be set from 5 to 720 minutes, 0 means no retry.
* `dst_info` - (Optional, List) Target information, single-node database use.
* `expect_run_time` - (Optional, String) Expected start time, when the value of RunMode is Timed, this value is required, such as: 2006-01-02 15:04:05.
* `job_mode` - (Optional, String) The enumeration values are liteMode and fullMode, corresponding to lite mode or normal mode respectively.
* `job_name` - (Optional, String) Sync job name.
* `options` - (Optional, List) Sync Task Options.
* `run_mode` - (Optional, String) Operation mode, such as: Immediate (indicates immediate operation, the default value is this value), Timed (indicates scheduled operation).
* `src_info` - (Optional, List) Source information, single-node database use.

The `columns` object of `tables` supports the following:

* `column_name` - (Optional, String) Column name. Note: This field may return null, indicating that no valid value can be obtained.
* `new_column_name` - (Optional, String) New column name. Note: This field may return null, indicating that no valid value can be obtained.

The `conflict_handle_option` object of `options` supports the following:

* `condition_column` - (Optional, String) Columns covered by the condition. Note: This field may return null, indicating that no valid value can be obtained.
* `condition_operator` - (Optional, String) Conditional Override Operation. Note: This field may return null, indicating that no valid value can be obtained.
* `condition_order_in_src_and_dst` - (Optional, String) Conditional Override Priority Processing. Note: This field may return null, indicating that no valid value can be obtained.

The `databases` object of `objects` supports the following:

* `db_mode` - (Optional, String) DB selection mode: All (for all objects under the current object), Partial (for some objects), when the Mode is Partial, this item is required. Note that synchronization of advanced objects does not depend on this value. Note: This field may return null, indicating that no valid value can be obtained.
* `db_name` - (Optional, String) The name of the library that needs to be migrated or synchronized. This item is required when the ObjectMode is Partial. Note: This field may return null, indicating that no valid value can be obtained.
* `event_mode` - (Optional, String) Event migration mode, all (for all objects under the current object), partial (partial objects). Note: This field may return null, indicating that no valid value can be obtained.
* `events` - (Optional, Set) When EventMode is partial, specify the name of the event to be migrated. Note: This field may return null, indicating that no valid value can be obtained.
* `function_mode` - (Optional, String) Select the mode to be synchronized, Partial is a part, all is an entire selection. Note: This field may return null, indicating that no valid value can be obtained.
* `functions` - (Optional, Set) Required when the FunctionMode value is Partial. Note: This field may return null, indicating that no valid value can be obtained.
* `new_db_name` - (Optional, String) The name of the library after migration or synchronization, which is the same as the source library by default. Note: This field may return null, indicating that no valid value can be obtained.
* `new_schema_name` - (Optional, String) Schema name after migration or synchronization. Note: This field may return null, indicating that no valid value can be obtained.
* `procedure_mode` - (Optional, String) Select the mode to be synchronized, Partial is part, All is the whole selection. Note: This field may return null, indicating that no valid value can be obtained.
* `procedures` - (Optional, Set) Required when the value of ProcedureMode is Partial. Note: This field may return null, indicating that no valid value can be obtained.
* `schema_mode` - (Optional, String) Schema selection mode, used by PostgreSQL and SQL Server sync links. Valid values: `All` (all objects under the current object), `Partial` (some objects). Note: This field may return null, indicating that no valid value can be obtained.
* `schema_name` - (Optional, String) Migrated or synchronized schemaNote: This field may return null, indicating that no valid value can be obtained.
* `table_mode` - (Optional, String) Table selection mode: All (for all objects under the current object), Partial (for some objects), this item is required when the DBMode is Partial. Note: This field may return null, indicating that no valid value can be obtained.
* `tables` - (Optional, List) A collection of table graph objects, when TableMode is Partial, this item needs to be filled in. Note: This field may return null, indicating that no valid value can be obtained.
* `trigger_mode` - (Optional, String) Trigger migration mode, all (for all objects under the current object), partial (partial objects). Note: This field may return null, indicating that no valid value can be obtained.
* `triggers` - (Optional, Set) When TriggerMode is partial, specify the name of the trigger to be migrated. Note: This field may return null, indicating that no valid value can be obtained.
* `view_mode` - (Optional, String) View selection mode: All is all view objects under the current object, Partial is part of the view objects. Note: This field may return null, indicating that no valid value can be obtained.
* `views` - (Optional, List) View object collection, when ViewMode is Partial, this item needs to be filled in. Note: This field may return null, indicating that no valid value can be obtained.

The `ddl_options` object of `options` supports the following:

* `ddl_object` - (Optional, String) Ddl type, such as Database, Table, View, Index, etc. Note: This field may return null, indicating that no valid value can be obtained.
* `ddl_value` - (Optional, Set) The specific value of ddl, the possible values for Database [Create,Drop,Alter].The possible values for Table [Create,Drop,Alter,Truncate,Rename].The possible values for View[Create,Drop].For the possible values of Index [Create, Drop]. Note: This field may return null, indicating that no valid value can be obtained.

The `dst_info` object supports the following:

* `account_mode` - (Optional, String) The account to which the resource belongs is empty or self (represents resources within this account), other (represents cross-account resources). Note: This field may return null, indicating that no valid value can be obtained.
* `account_role` - (Optional, String) The role during cross-account synchronization, only [a-zA-Z0-9-_]+ is allowed, if it is a cross-account instance, this field is required. Note: This field may return null, indicating that no valid value can be obtained.
* `account` - (Optional, String) The account to which the instance belongs. This field is required if it is a cross-account instance. Note: This field may return null, indicating that no valid value can be obtained.
* `ccn_id` - (Optional, String) Cloud networking ID, which is required for the cloud networking access type. Note: This field may return null, indicating that no valid value can be obtained.
* `cvm_instance_id` - (Optional, String) CVM instance short ID, which is the same as the instance ID displayed on the cloud server console page. If it is a self-built instance of CVM, this field needs to be passed. Note: This field may return null, indicating that no valid value can be obtained.
* `database_net_env` - (Optional, String) The network environment to which the database belongs. It is required when AccessType is Cloud Network (CCN). `UserIDC` represents the user IDC. `TencentVPC` represents Tencent Cloud VPC. Note: This field may return null, indicating that no valid value can be obtained.
* `db_kernel` - (Optional, String) Database kernel type, used to distinguish different kernels in tdsql: percona, mariadb, mysql. Note: This field may return null, indicating that no valid value can be obtained.
* `db_name` - (Optional, String) Database name, when the database is cdwpg, it needs to be provided. Note: This field may return null, indicating that no valid value can be obtained.
* `encrypt_conn` - (Optional, String) Whether to use encrypted transmission, UnEncrypted means not to use encrypted transmission, Encrypted means to use encrypted transmission, the default is UnEncrypted. Note: This field may return null, indicating that no valid value can be obtained.
* `engine_version` - (Optional, String) Database version, valid only when the instance is an RDS instance, ignored by other instances, the format is: 5.6 or 5.7, the default is 5.6. Note: This field may return null, indicating that no valid value can be obtained.
* `instance_id` - (Optional, String) Database instance id. Note: This field may return null, indicating that no valid value can be obtained.
* `ip` - (Optional, String) The IP address of the instance, which is required when the access type is non-cdb. Note: This field may return null, indicating that no valid value can be obtained.
* `password` - (Optional, String) Password, required for instances that require username and password authentication for access. Note: This field may return null, indicating that no valid value can be obtained.
* `port` - (Optional, Int) Instance port, this item is required when the access type is non-cdb. Note: This field may return null, indicating that no valid value can be obtained.
* `region` - (Optional, String) The english name of region. Note: This field may return null, indicating that no valid value can be obtained.
* `role_external_id` - (Optional, String) External role id. Note: This field may return null, indicating that no valid value can be obtained.
* `role` - (Optional, String) The node type of tdsql mysql version, the enumeration value is proxy, set. Note: This field may return null, indicating that no valid value can be obtained.
* `subnet_id` - (Optional, String) The subnet ID under the private network, this item is required for the private network, leased line, and VPN access methods. Note: This field may return null, indicating that no valid value can be obtained.
* `supplier` - (Optional, String) Cloud vendor type, when the instance is an RDS instance, fill in aliyun, in other cases fill in others, the default is others. Note: This field may return null, indicating that no valid value can be obtained.
* `tmp_secret_id` - (Optional, String) Temporary key Id, required if it is a cross-account instance. Note: This field may return null, indicating that no valid value can be obtained.
* `tmp_secret_key` - (Optional, String) Temporary key Key, required if it is a cross-account instance. Note: This field may return null, indicating that no valid value can be obtained.
* `tmp_token` - (Optional, String) Temporary Token, required if it is a cross-account instance. Note: This field may return null, indicating that no valid value can be obtained.
* `uniq_dcg_id` - (Optional, String) Leased line gateway ID, which is required for the leased line access type. Note: This field may return null, indicating that no valid value can be obtained.
* `uniq_vpn_gw_id` - (Optional, String) VPN gateway ID, which is required for the VPN access type. Note: This field may return null, indicating that no valid value can be obtained.
* `user` - (Optional, String) Username, required for instances that require username and password authentication for access. Note: This field may return null, indicating that no valid value can be obtained.
* `vpc_id` - (Optional, String) Private network ID, which is required for access methods of private network, leased line, and VPN. Note: This field may return null, indicating that no valid value can be obtained.

The `objects` object supports the following:

* `advanced_objects` - (Optional, Set) For advanced object types, such as function and procedure, when an advanced object needs to be synchronized, the initialization type must include the structure initialization type, that is, the value of the Options.InitType field is Structure or Full. Note: This field may return null, indicating that no valid value can be obtained.
* `databases` - (Optional, List) Synchronization object, not null when Mode is Partial. Note: This field may return null, indicating that no valid value can be obtained.
* `mode` - (Optional, String) Migration object type Partial (partial object). Note: This field may return null, indicating that no valid value can be obtained.
* `online_ddl` - (Optional, List, **Deprecated**) It has been deprecated from version 1.83.20. OnlineDDL type. Note: This field may return null, indicating that no valid value can be obtained.

The `online_ddl` object of `objects` supports the following:

* `status` - (Optional, String, **Deprecated**) It has been deprecated from version 1.83.20. Status: ON - Enabled, OFF - Disabled.

The `options` object supports the following:

* `add_additional_column` - (Optional, Bool) Whether to add additional columns. Note: This field may return null, indicating that no valid value can be obtained.
* `conflict_handle_option` - (Optional, List) Detailed options for conflict handling, such as conditional rows and conditional actions in conditional overrides. Note: This field may return null, indicating that no valid value can be obtained.
* `conflict_handle_type` - (Optional, String) Conflict handling options, ReportError (error report, the default value), Ignore (ignore), Cover (cover), ConditionCover (condition coverage). Note: This field may return null, indicating that no valid value can be obtained.
* `ddl_options` - (Optional, List) DDL synchronization options, specifically describe which DDLs to synchronize. Note: This field may return null, indicating that no valid value can be obtained.
* `deal_of_exist_same_table` - (Optional, String) The processing of the table with the same name, ReportErrorAfterCheck (pre-check and report error, default), InitializeAfterDelete (delete and re-initialize), ExecuteAfterIgnore (ignore and continue to execute). Note: This field may return null, indicating that no valid value can be obtained.
* `init_type` - (Optional, String) Synchronous initialization options, Data (full data initialization), Structure (structure initialization), Full (full data and structure initialization, default), None (incremental only). Note: This field may return null, indicating that no valid value can be obtained.
* `op_types` - (Optional, Set) DML and DDL options to be synchronized, Insert (insert operation), Update (update operation), Delete (delete operation), DDL (structure synchronization), leave blank (not selected), PartialDDL (custom, work with DdlOptions). Note: This field may return null, indicating that no valid value can be obtained.
* `rate_limit_option` - (Optional, List) Task speed limit information
Note: This field may return null, indicating that no valid values can be obtained.

The `rate_limit_option` object of `options` supports the following:

* `current_dump_rps` - (Required, Int) The full export Rps currently in effect. The value of this field can be adjusted when configuring the task. Note: If it is not set or set to 0, it means the current value is maintained. The maximum value is 50,000,000.Note: This field may return null, indicating that no valid values can be obtained.
* `current_dump_thread` - (Required, Int) The number of full export threads currently in effect. The value of this field can be adjusted when configuring the task. Note: If it is not set or set to 0, it means the current value is maintained. The maximum value is 16.Note: This field may return null, indicating that no valid values can be obtained.
* `current_load_rps` - (Required, Int) The full import Rps currently in effect. The value of this field can be adjusted when configuring the task. Note: If it is not set or set to 0, it means the current value is maintained. The maximum value is 50,000,000.Note: This field may return null, indicating that no valid values can be obtained.
* `current_load_thread` - (Required, Int) The number of full import threads currently in effect. The value of this field can be adjusted when configuring the task. Note: If it is not set or set to 0, it means the current value is maintained. The maximum value is 16.Note: This field may return null, indicating that no valid values can be obtained.
* `current_sinker_thread` - (Required, Int) The number of incremental import threads currently in effect. The value of this field can be adjusted when configuring the task. Note: If it is not set or set to 0, it means the current value is maintained. The maximum value is 128.Note: This field may return null, indicating that no valid values can be obtained.
* `default_dump_rps` - (Required, Int) The default full export Rps. This field is only meaningful in the output parameter.Note: This field may return null, indicating that no valid values can be obtained.
* `default_dump_thread` - (Required, Int) The default number of full export threads. This field is only meaningful in the output parameter.Note: This field may return null, indicating that no valid values can be obtained.
* `default_load_rps` - (Required, Int) The default full import Rps. This field is only meaningful in the output parameter.Note: This field may return null, indicating that no valid values can be obtained.
* `default_load_thread` - (Required, Int) The default number of full import threads. This field is only meaningful in the output parameter.Note: This field may return null, indicating that no valid values can be obtained.
* `default_sinker_thread` - (Required, Int) The default number of incremental import threads. This field is only meaningful in the output parameter.
Note: This field may return null, indicating that no valid values can be obtained.
* `has_user_set_rate_limit` - (Required, String) enum:"no"/"yes", no: the user has not set a speed limit; yes: a speed limit has been set. This field is only meaningful in the output parameter.Note: This field may return null, indicating that no valid values can be obtained.

The `src_info` object supports the following:

* `account_mode` - (Optional, String) The account to which the resource belongs is empty or self (represents resources within this account), other (represents cross-account resources). Note: This field may return null, indicating that no valid value can be obtained.
* `account_role` - (Optional, String) The role during cross-account synchronization, only [a-zA-Z0-9-_]+ is allowed, if it is a cross-account instance, this field is required. Note: This field may return null, indicating that no valid value can be obtained.
* `account` - (Optional, String) The account to which the instance belongs. This field is required if it is a cross-account instance. Note: This field may return null, indicating that no valid value can be obtained.
* `ccn_id` - (Optional, String) Cloud networking ID, which is required for the cloud networking access type. Note: This field may return null, indicating that no valid value can be obtained.
* `cvm_instance_id` - (Optional, String) CVM instance short ID, which is the same as the instance ID displayed on the cloud server console page. If it is a self-built instance of CVM, this field needs to be passed. Note: This field may return null, indicating that no valid value can be obtained.
* `database_net_env` - (Optional, String) The network environment to which the database belongs. It is required when AccessType is Cloud Network (CCN). `UserIDC` represents the user IDC. `TencentVPC` represents Tencent Cloud VPC. Note: This field may return null, indicating that no valid value can be obtained.
* `db_kernel` - (Optional, String) Database kernel type, used to distinguish different kernels in tdsql: percona, mariadb, mysql. Note: This field may return null, indicating that no valid value can be obtained.
* `db_name` - (Optional, String) Database name, when the database is cdwpg, it needs to be provided. Note: This field may return null, indicating that no valid value can be obtained.
* `encrypt_conn` - (Optional, String) Whether to use encrypted transmission, UnEncrypted means not to use encrypted transmission, Encrypted means to use encrypted transmission, the default is UnEncrypted. Note: This field may return null, indicating that no valid value can be obtained.
* `engine_version` - (Optional, String) Database version, valid only when the instance is an RDS instance, ignored by other instances, the format is: 5.6 or 5.7, the default is 5.6. Note: This field may return null, indicating that no valid value can be obtained.
* `instance_id` - (Optional, String) Database instance id. Note: This field may return null, indicating that no valid value can be obtained.
* `ip` - (Optional, String) The IP address of the instance, which is required when the access type is non-cdb. Note: This field may return null, indicating that no valid value can be obtained.
* `password` - (Optional, String) Password, required for instances that require username and password authentication for access. Note: This field may return null, indicating that no valid value can be obtained.
* `port` - (Optional, Int) Instance port, this item is required when the access type is non-cdb. Note: This field may return null, indicating that no valid value can be obtained.
* `region` - (Optional, String) The english name of region. Note: This field may return null, indicating that no valid value can be obtained.
* `role_external_id` - (Optional, String) External role id. Note: This field may return null, indicating that no valid value can be obtained.
* `role` - (Optional, String) The node type of tdsql mysql version, the enumeration value is proxy, set. Note: This field may return null, indicating that no valid value can be obtained.
* `subnet_id` - (Optional, String) The subnet ID under the private network, this item is required for the private network, leased line, and VPN access methods. Note: This field may return null, indicating that no valid value can be obtained.
* `supplier` - (Optional, String) Cloud vendor type, when the instance is an RDS instance, fill in aliyun, in other cases fill in others, the default is others. Note: This field may return null, indicating that no valid value can be obtained.
* `tmp_secret_id` - (Optional, String) Temporary key Id, required if it is a cross-account instance. Note: This field may return null, indicating that no valid value can be obtained.
* `tmp_secret_key` - (Optional, String) Temporary key Key, required if it is a cross-account instance. Note: This field may return null, indicating that no valid value can be obtained.
* `tmp_token` - (Optional, String) Temporary Token, required if it is a cross-account instance. Note: This field may return null, indicating that no valid value can be obtained.
* `uniq_dcg_id` - (Optional, String) Leased line gateway ID, which is required for the leased line access type. Note: This field may return null, indicating that no valid value can be obtained.
* `uniq_vpn_gw_id` - (Optional, String) VPN gateway ID, which is required for the VPN access type. Note: This field may return null, indicating that no valid value can be obtained.
* `user` - (Optional, String) Username, required for instances that require username and password authentication for access. Note: This field may return null, indicating that no valid value can be obtained.
* `vpc_id` - (Optional, String) Private network ID, which is required for access methods of private network, leased line, and VPN. Note: This field may return null, indicating that no valid value can be obtained.

The `tables` object of `databases` supports the following:

* `column_mode` - (Optional, String) Whether to synchronize all columns in the table, All: all columns under the current table, Partial: part of the columns under the current table, fill in the detailed column information through the columns field. Note: This field may return null, indicating that no valid value can be obtained.
* `columns` - (Optional, List) Synchronized column information, required when ColumnMode is Partial. Note: This field may return null, indicating that no valid value can be obtained.
* `filter_condition` - (Optional, String) Filter condition. Note: This field may return null, indicating that no valid value can be obtained.
* `new_table_name` - (Optional, String) New table name. Note: This field may return null, indicating that no valid value can be obtained.
* `table_edit_mode` - (Optional, String) Edit table type, rename (table mapping), pt (synchronize attached tables). Note: This field may return null, indicating that no valid value can be obtained.
* `table_name` - (Optional, String) Table name. Note: This field may return null, indicating that no valid value can be obtained.
* `tmp_tables` - (Optional, Set) Temporary table generated by pt-osc/gh-ost. Note: This field may return null, indicating that no valid value can be obtained.

The `views` object of `databases` supports the following:

* `new_view_name` - (Optional, String) New view name. Note: This field may return null, indicating that no valid value can be obtained.
* `view_name` - (Optional, String) View name. Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

DTS sync config can be imported using the id, e.g.

```
terraform import tencentcloud_dts_sync_config.example sync-muu9ez38
```

