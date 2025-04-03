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

* `dst_info` - (Required, List) Target database information.
* `migrate_option` - (Required, List) Migration job configuration options, used to describe how the task performs migration.
* `run_mode` - (Required, String) Running mode. Valid values: immediate, timed.
* `service_id` - (Required, String) Migrate service Id from `tencentcloud_dts_migrate_service`.
* `src_info` - (Required, List) Source instance information.
* `auto_retry_time_range_minutes` - (Optional, Int) The automatic retry time period can be set from 5 to 720 minutes, with 0 indicating no retry.
* `expect_run_time` - (Optional, String) Expected start time in the format of `2006-01-02 15:04:05`, which is required if RunMode is timed.

The `consistency` object of `migrate_option` supports the following:

* `mode` - (Optional, String) Data consistency check type. Valid values: full, noCheck, notConfigured.

The `database_table` object of `migrate_option` supports the following:

* `object_mode` - (Required, String) Migration object type. Valid values: all, partial.
* `advanced_objects` - (Optional, Set) Advanced object types, such as trigger, function, procedure, event. Note: If you want to migrate and synchronize advanced objects, the corresponding advanced object type should be included in this configuration.
* `databases` - (Optional, List) Migration object, which is required if ObjectMode is partial.

The `databases` object of `database_table` supports the following:

* `db_mode` - (Optional, String) Database selection mode, which is required if ObjectMode is partial. Valid values: all, partial.
* `db_name` - (Optional, String) Name of the database to be migrated or synced, which is required if ObjectMode is partial.
* `event_mode` - (Optional, String) Sync mode. Valid values: partial, all.
* `events` - (Optional, Set) This parameter is required if EventMode is partial.
* `function_mode` - (Optional, String) Sync mode. Valid values: partial, all.
* `functions` - (Optional, Set) This parameter is required if FunctionMode is partial.
* `new_db_name` - (Optional, String) Name of the database after migration or sync, which is the same as the source database name by default.
* `new_schema_name` - (Optional, String) Name of the schema after migration or sync.
* `procedure_mode` - (Optional, String) Sync mode. Valid values: partial, all.
* `procedures` - (Optional, Set) This parameter is required if ProcedureMode is partial.
* `role_mode` - (Optional, String) Role selection mode, which is exclusive to PostgreSQL. Valid values: all, partial.
* `roles` - (Optional, List) Role, which is exclusive to PostgreSQL and required if RoleMode is partial.
* `schema_mode` - (Optional, String) Schema selection mode. Valid values: all, partial.
* `schema_name` - (Optional, String) The schema to be migrated or synced.
* `table_mode` - (Optional, String) Table selection mode, which is required if DBMode is partial. Valid values: all, partial.
* `tables` - (Optional, List) The set of table objects, which is required if TableMode is partial.
* `trigger_mode` - (Optional, String) Sync mode. Valid values: partial, all.
* `triggers` - (Optional, Set) This parameter is required if TriggerMode is partial.
* `view_mode` - (Optional, String) View selection mode. Valid values: all, partial.
* `views` - (Optional, List) The set of view objects, which is required if ViewMode is partial.

The `dst_info` object supports the following:

* `access_type` - (Required, String) Instances network access type. Valid values: extranet (public network); ipv6 (public IPv6); cvm (self-build on CVM); dcg (Direct Connect); vpncloud (VPN access); cdb (database); ccn (CCN); intranet (intranet); vpc (VPC). Note that the valid values are subject to the current link.
* `database_type` - (Required, String) Database type, such as mysql, redis, mongodb, postgresql, mariadb, and percona.
* `info` - (Required, List) Database information.
* `node_type` - (Required, String) Node type, empty or simple indicates a general node, cluster indicates a cluster node; for mongo services, valid values: replicaset (mongodb replica set), standalone (mongodb single node), cluster (mongodb cluster); for redis instances, valid values: empty or simple (single node), cluster (cluster), cluster-cache (cache cluster), cluster-proxy (proxy cluster).
* `region` - (Required, String) Instance region.
* `extra_attr` - (Optional, List) For MongoDB, you can define the following parameters: ['AuthDatabase':'admin','AuthFlag': '1', 'AuthMechanism':'SCRAM-SHA-1'].
* `supplier` - (Optional, String) Instance service provider, such as `aliyun` and `others`.

The `extra_attr` object of `dst_info` supports the following:

* `key` - (Optional, String) Option key.
* `value` - (Optional, String) Option value.

The `extra_attr` object of `migrate_option` supports the following:

* `key` - (Optional, String) Option key.
* `value` - (Optional, String) Option value.

The `extra_attr` object of `src_info` supports the following:

* `key` - (Optional, String) Option key.
* `value` - (Optional, String) Option value.

The `info` object of `dst_info` supports the following:

* `account_mode` - (Optional, String) The account to which the resource belongs. Valid values: empty or self (the current account); other (another account).
* `account_role` - (Optional, String) The role used for cross-account migration, which can contain [a-zA-Z0-9-_]+.
* `account` - (Optional, String) Instance account.
* `ccn_gw_id` - (Optional, String) CCN instance ID such as ccn-afp6kltc.
* `cvm_instance_id` - (Optional, String) Short CVM instance ID in the format of ins-olgl39y8, which is required if the access type is cvm. It is the same as the instance ID displayed in the CVM console.
* `db_kernel` - (Optional, String) Kernel version, such as the different kernel versions of MariaDB.
* `engine_version` - (Optional, String) Database version in the format of 5.6 or 5.7, which takes effect only if the instance is an RDS instance. Default value: 5.6.
* `host` - (Optional, String) Instance IP address, which is required for the following access types: public network, Direct Connect, VPN, CCN, intranet, and VPC.
* `instance_id` - (Optional, String) Database instance ID in the format of cdb-powiqx8q, which is required if the access type is cdb.
* `password` - (Optional, String) Instance password.
* `port` - (Optional, Int) Instance port, which is required for the following access types: public network, self-build on CVM, Direct Connect, VPN, CCN, intranet, and VPC.
* `role` - (Optional, String) Node role in a distributed database, such as the mongos node in MongoDB.
* `subnet_id` - (Optional, String) ID of the subnet in the VPC in the format of subnet-3paxmkdz, which is required if the access type is vpc, vpncloud, ccn, or dcg.
* `tmp_secret_id` - (Optional, String) Temporary SecretId, you can obtain the temporary key by GetFederationToken.
* `tmp_secret_key` - (Optional, String) Temporary SecretKey, you can obtain the temporary key by GetFederationToken.
* `tmp_token` - (Optional, String) Temporary token, you can obtain the temporary key by GetFederationToken.
* `uniq_dcg_id` - (Optional, String) Direct Connect gateway ID in the format of dcg-0rxtqqxb, which is required if the access type is dcg.
* `uniq_vpn_gw_id` - (Optional, String) VPN gateway ID in the format of vpngw-9ghexg7q, which is required if the access type is vpncloud.
* `user` - (Optional, String) Instance username.
* `vpc_id` - (Optional, String) VPC ID in the format of vpc-92jblxto, which is required if the access type is vpc, vpncloud, ccn, or dcg.

The `info` object of `src_info` supports the following:

* `account_mode` - (Optional, String) The account to which the resource belongs. Valid values: empty or self (the current account); other (another account).
* `account_role` - (Optional, String) The role used for cross-account migration, which can contain [a-zA-Z0-9-_]+.
* `account` - (Optional, String) Instance account.
* `ccn_gw_id` - (Optional, String) CCN instance ID such as ccn-afp6kltc.
* `cvm_instance_id` - (Optional, String) Short CVM instance ID in the format of ins-olgl39y8, which is required if the access type is cvm. It is the same as the instance ID displayed in the CVM console.
* `db_kernel` - (Optional, String) Kernel version, such as the different kernel versions of MariaDB.
* `engine_version` - (Optional, String) Database version in the format of 5.6 or 5.7, which takes effect only if the instance is an RDS instance. Default value: 5.6.
* `host` - (Optional, String) Instance IP address, which is required for the following access types: public network, Direct Connect, VPN, CCN, intranet, and VPC.
* `instance_id` - (Optional, String) Database instance ID in the format of cdb-powiqx8q, which is required if the access type is cdb.
* `password` - (Optional, String) Instance password.
* `port` - (Optional, Int) Instance port, which is required for the following access types: public network, self-build on CVM, Direct Connect, VPN, CCN, intranet, and VPC.
* `role` - (Optional, String) Node role in a distributed database, such as the mongos node in MongoDB.
* `subnet_id` - (Optional, String) ID of the subnet in the VPC in the format of subnet-3paxmkdz, which is required if the access type is vpc, vpncloud, ccn, or dcg.
* `tmp_secret_id` - (Optional, String) Temporary SecretId, you can obtain the temporary key by GetFederationToken.
* `tmp_secret_key` - (Optional, String) Temporary SecretKey, you can obtain the temporary key by GetFederationToken.
* `tmp_token` - (Optional, String) Temporary token, you can obtain the temporary key by GetFederationToken.
* `uniq_dcg_id` - (Optional, String) Direct Connect gateway ID in the format of dcg-0rxtqqxb, which is required if the access type is dcg.
* `uniq_vpn_gw_id` - (Optional, String) VPN gateway ID in the format of vpngw-9ghexg7q, which is required if the access type is vpncloud.
* `user` - (Optional, String) Instance username.
* `vpc_id` - (Optional, String) VPC ID in the format of vpc-92jblxto, which is required if the access type is vpc, vpncloud, ccn, or dcg.

The `migrate_option` object supports the following:

* `database_table` - (Required, List) Migration object option, you need to tell the migration service which library table objects to migrate.
* `consistency` - (Optional, List) Data consistency check option. Data consistency check is disabled by default.
* `extra_attr` - (Optional, List) Additional information. You can set additional parameters for certain database types.
* `is_dst_read_only` - (Optional, Bool) Whether to set the target database to read-only during migration, which takes effect only for MySQL databases. Valid values: true, false. Default value: false.
* `is_migrate_account` - (Optional, Bool) Whether to migrate accounts.
* `is_override_root` - (Optional, Bool) Whether to use the Root account in the source database to overwrite that in the target database. Valid values: false, true. For database/table or structural migration, you should specify false. Note that this parameter takes effect only for OldDTS.
* `migrate_type` - (Optional, String) Migration type. Valid values: full, structure, fullAndIncrement. Default value: fullAndIncrement.

The `roles` object of `databases` supports the following:

* `new_role_name` - (Optional, String) Role name after migration.
* `role_name` - (Optional, String) Role name.

The `src_info` object supports the following:

* `access_type` - (Required, String) Instances network access type. Valid values: extranet (public network); ipv6 (public IPv6); cvm (self-build on CVM); dcg (Direct Connect); vpncloud (VPN access); cdb (database); ccn (CCN); intranet (intranet); vpc (VPC). Note that the valid values are subject to the current link.
* `database_type` - (Required, String) Database type, such as mysql, redis, mongodb, postgresql, mariadb, and percona.
* `info` - (Required, List) Database information.
* `node_type` - (Required, String) Node type, empty or simple indicates a general node, cluster indicates a cluster node; for mongo services, valid values: replicaset (mongodb replica set), standalone (mongodb single node), cluster (mongodb cluster); for redis instances, valid values: empty or simple (single node), cluster (cluster), cluster-cache (cache cluster), cluster-proxy (proxy cluster).
* `region` - (Required, String) Instance region.
* `extra_attr` - (Optional, List) For MongoDB, you can define the following parameters: ['AuthDatabase':'admin', 'AuthFlag': '1', 'AuthMechanism':'SCRAM-SHA-1'].
* `supplier` - (Optional, String) Instance service provider, such as `aliyun` and `others`.

The `tables` object of `databases` supports the following:

* `new_table_name` - (Optional, String) New name of the migrated table. This parameter is required when TableEditMode is rename. It is mutually exclusive with TmpTables..
* `table_edit_mode` - (Optional, String) Table editing type. Valid values: rename (table mapping); pt (additional table sync).
* `table_name` - (Optional, String) Name of the migrated table, which is case-sensitive.
* `tmp_tables` - (Optional, Set) The temp tables to be migrated. This parameter is mutually exclusive with NewTableName. It is valid only when the configured migration objects are table-level ones and TableEditMode is pt. To migrate temp tables generated when pt-osc or other tools are used during the migration process, you must configure this parameter first. For example, if you want to perform the pt-osc operation on a table named 't1', configure this parameter as ['_t1_new','_t1_old']; to perform the gh-ost operation on t1, configure it as ['_t1_ghc','_t1_gho','_t1_del']. Temp tables generated by pt-osc and gh-ost operations can be configured at the same time.

The `views` object of `databases` supports the following:

* `new_view_name` - (Optional, String) View name after migration.
* `view_name` - (Optional, String) View name.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `status` - Task status. Valid values: created(Created), checking (Checking), checkPass (Check passed), checkNotPass (Check not passed), readyRun (Ready for running), running (Running), readyComplete (Preparation completed), success (Successful), failed (Failed), stopping (Stopping), completing (Completing), pausing (Pausing), manualPaused (Paused).


## Import

DTS migrate job can be imported using the id, e.g.

```
terraform import tencentcloud_dts_migrate_job.example dts-iy98oxba
```

