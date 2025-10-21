---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_sync_config"
sidebar_current: "docs-tencentcloud-resource-dts_sync_config"
description: |-
  Provides a resource to create a dts sync_config
---

# tencentcloud_dts_sync_config

Provides a resource to create a dts sync_config

## Example Usage

### Sync mysql database to cynosdb through cdb access type

```hcl
resource "tencentcloud_cynosdb_cluster" "foo" {
  available_zone               = var.availability_zone
  vpc_id                       = local.vpc_id
  subnet_id                    = local.subnet_id
  db_type                      = "MYSQL"
  db_version                   = "5.7"
  storage_limit                = 1000
  cluster_name                 = "tf-cynosdb-mysql-sync-dst"
  password                     = "*"
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

resource "tencentcloud_dts_sync_job" "sync_job" {
  pay_mode          = "PostPay"
  src_database_type = "mysql"
  src_region        = "ap-guangzhou"
  dst_database_type = "cynosdbmysql"
  dst_region        = "ap-guangzhou"
  tags {
    tag_key   = "aaa"
    tag_value = "bbb"
  }
  auto_renew     = 0
  instance_class = "micro"
}

resource "tencentcloud_dts_sync_config" "sync_config" {
  job_id          = tencentcloud_dts_sync_job.sync_job.job_id
  src_access_type = "cdb"
  dst_access_type = "cdb"

  job_name = "tf_test_sync_config"
  job_mode = "liteMode"
  run_mode = "Immediate"

  objects {
    mode = "Partial"
    databases {
      db_name     = "tf_ci_test"
      new_db_name = "tf_ci_test_new"
      db_mode     = "Partial"
      table_mode  = "All"
      tables {
        table_name     = "test"
        new_table_name = "test_new"
      }
    }
  }
  src_info {
    region      = "ap-guangzhou"
    instance_id = "cdb-fitq5t9h"
    user        = "your_user_name"
    password    = "*"
    db_name     = "tf_ci_test"
    vpc_id      = local.vpc_id
    subnet_id   = local.subnet_id
  }
  dst_info {
    region      = "ap-guangzhou"
    instance_id = tencentcloud_cynosdb_cluster.foo.id
    user        = "root"
    password    = "*"
    db_name     = "tf_ci_test_new"
    vpc_id      = local.vpc_id
    subnet_id   = local.subnet_id
  }
  auto_retry_time_range_minutes = 0
}
```

### Sync mysql database using CCN to route from ap-shanghai to ap-guangzhou

```hcl
locals {
  vpc_id_sh    = "vpc-evtcyb3g"
  subnet_id_sh = "subnet-1t83cxkp"
  src_ip       = data.tencentcloud_mysql_instance.src_mysql.instance_list.0.intranet_ip
  src_port     = data.tencentcloud_mysql_instance.src_mysql.instance_list.0.intranet_port
  ccn_id       = data.tencentcloud_ccn_instances.ccns.instance_list.0.ccn_id
  dst_mysql_id = data.tencentcloud_mysql_instance.dst_mysql.instance_list.0.mysql_id
}

variable "src_az_sh" {
  default = "ap-shanghai"
}

variable "dst_az_gz" {
  default = "ap-guangzhou"
}

data "tencentcloud_dts_sync_jobs" "sync_jobs" {
  job_name = "keep_sync_config_ccn_2_cdb"
}

data "tencentcloud_ccn_instances" "ccns" {
  name = "keep-ccn-dts-sh"
}

data "tencentcloud_mysql_instance" "src_mysql" {
  instance_name = "your_user_name_mysql_src"
}

data "tencentcloud_mysql_instance" "dst_mysql" {
  instance_name = "your_user_name_mysql_src"
}

resource "tencentcloud_dts_sync_config" "sync_config" {
  job_id          = data.tencentcloud_dts_sync_jobs.sync_jobs.list.0.job_id
  src_access_type = "ccn"
  dst_access_type = "cdb"

  job_mode = "liteMode"
  run_mode = "Immediate"

  objects {
    mode = "Partial"
    databases {
      db_name     = "tf_ci_test"
      new_db_name = "tf_ci_test_new"
      db_mode     = "Partial"
      table_mode  = "All"
      tables {
        table_name     = "test"
        new_table_name = "test_new"
      }
    }
  }
  src_info { // shanghai to guangzhou via ccn
    region           = var.src_az_sh
    user             = "your_user_name"
    password         = "your_pass_word"
    ip               = local.src_ip
    port             = local.src_port
    vpc_id           = local.vpc_id_sh
    subnet_id        = local.subnet_id_sh
    ccn_id           = local.ccn_id
    database_net_env = "TencentVPC"
  }
  dst_info {
    region      = var.dst_az_gz
    instance_id = local.dst_mysql_id
    user        = "your_user_name"
    password    = "your_pass_word"
  }
  auto_retry_time_range_minutes = 0
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
* `online_ddl` - (Optional, List) OnlineDDL type. Note: This field may return null, indicating that no valid value can be obtained.

The `online_ddl` object of `objects` supports the following:

* `status` - (Optional, String) status.

The `options` object supports the following:

* `add_additional_column` - (Optional, Bool) Whether to add additional columns. Note: This field may return null, indicating that no valid value can be obtained.
* `conflict_handle_option` - (Optional, List) Detailed options for conflict handling, such as conditional rows and conditional actions in conditional overrides. Note: This field may return null, indicating that no valid value can be obtained.
* `conflict_handle_type` - (Optional, String) Conflict handling options, ReportError (error report, the default value), Ignore (ignore), Cover (cover), ConditionCover (condition coverage). Note: This field may return null, indicating that no valid value can be obtained.
* `ddl_options` - (Optional, List) DDL synchronization options, specifically describe which DDLs to synchronize. Note: This field may return null, indicating that no valid value can be obtained.
* `deal_of_exist_same_table` - (Optional, String) The processing of the table with the same name, ReportErrorAfterCheck (pre-check and report error, default), InitializeAfterDelete (delete and re-initialize), ExecuteAfterIgnore (ignore and continue to execute). Note: This field may return null, indicating that no valid value can be obtained.
* `init_type` - (Optional, String) Synchronous initialization options, Data (full data initialization), Structure (structure initialization), Full (full data and structure initialization, default), None (incremental only). Note: This field may return null, indicating that no valid value can be obtained.
* `op_types` - (Optional, Set) DML and DDL options to be synchronized, Insert (insert operation), Update (update operation), Delete (delete operation), DDL (structure synchronization), leave blank (not selected), PartialDDL (custom, work with DdlOptions). Note: This field may return null, indicating that no valid value can be obtained.

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

* `filter_condition` - (Optional, String) Filter condition. Note: This field may return null, indicating that no valid value can be obtained.
* `new_table_name` - (Optional, String) New table name. Note: This field may return null, indicating that no valid value can be obtained.
* `table_name` - (Optional, String) Table name. Note: This field may return null, indicating that no valid value can be obtained.

The `views` object of `databases` supports the following:

* `new_view_name` - (Optional, String) New view name. Note: This field may return null, indicating that no valid value can be obtained.
* `view_name` - (Optional, String) View name. Note: This field may return null, indicating that no valid value can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dts sync_config can be imported using the id, e.g.

```
terraform import tencentcloud_dts_sync_config.sync_config sync_config_id
```

