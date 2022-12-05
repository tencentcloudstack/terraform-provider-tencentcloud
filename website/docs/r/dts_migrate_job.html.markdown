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
resource "tencentcloud_dts_migrate_job" "migrate_job" {
  src_database_type = "mysql"
  dst_database_type = "cynosdbmysql"
  src_region        = "ap-guangzhou"
  dst_region        = "ap-guangzhou"
  instance_class    = "small"
  job_name          = "tf_test_migration_job"
  tags {
    tag_key   = "aaa"
    tag_value = "bbb"
  }
}
```

## Argument Reference

The following arguments are supported:

* `dst_database_type` - (Required, String) destination database type, optional value is mysql/redis/percona/mongodb/postgresql/sqlserver/mariadb.
* `dst_region` - (Required, String) destination region.
* `instance_class` - (Required, String) instance class, optional value is small/medium/large/xlarge/2xlarge.
* `src_database_type` - (Required, String) source database type, optional value is mysql/redis/percona/mongodb/postgresql/sqlserver/mariadb.
* `src_region` - (Required, String) source region.
* `job_name` - (Optional, String) job name.
* `tags` - (Optional, List) tags.

The `tags` object supports the following:

* `tag_key` - (Optional, String) tag key.
* `tag_value` - (Optional, String) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `dst_info` - destination info.
  * `access_type` - access type.
  * `database_type` - database type.
  * `extra_attr` - extra attributes.
    * `key` - key.
    * `value` - value.
  * `info` - databse info list.
    * `account_mode` - account mode.
    * `account_role` - account role.
    * `account` - account.
    * `ccn_gw_id` - ccn gateway id.
    * `cvm_instance_id` - cvm instance id.
    * `db_kernel` - database kernel.
    * `engine_version` - engine version.
    * `host` - host.
    * `instance_id` - instance id.
    * `password` - password.
    * `port` - port.
    * `role` - node role.
    * `subnet_id` - subnet id.
    * `tmp_secret_id` - temporary secret id.
    * `tmp_secret_key` - temporary secret key.
    * `tmp_token` - temporary token.
    * `uniq_vpn_gw_id` - vpn gateway id.
    * `user` - user.
    * `vpc_id` - vpc id.
  * `node_type` - node type.
  * `region` - region.
  * `supplier` - supplier.
* `expect_run_time` - expected run time, such as 2006-01-02 15:04:05.
* `job_id` - job id.
* `migrate_option` - migrate option.
  * `consistency` - consistency option.
    * `mode` - mode, optional value is full/noCheck/notConfigure.
  * `database_table` - database table.
    * `databases` - database list.
      * `d_b_mode` - database mode.
      * `db_name` - database name.
      * `event_mode` - event mode.
      * `events` - event list.
      * `function_mode` - function mode.
      * `functions` - function list.
      * `new_db_name` - new database name.
      * `new_schema_name` - new schema name.
      * `procedure_mode` - procedure mode.
      * `procedures` - procedure list.
      * `role_mode` - role mode.
      * `roles` - role list.
        * `new_role_name` - new role name.
        * `role_name` - role name.
      * `schema_mode` - schema mode.
      * `schema_name` - schema name.
      * `table_mode` - table mode.
      * `tables` - table list.
        * `new_table_name` - new table name.
        * `table_edit_mode` - table edit mode.
        * `table_name` - table name.
        * `tmp_tables` - temporary tables.
      * `trigger_mode` - trigger mode.
      * `triggers` - trigger list.
      * `view_mode` - view mode.
      * `views` - views.
        * `new_view_name` - new view name.
        * `view_name` - view name.
    * `object_mode` - object mode.
  * `extra_attr` - extra attributes.
    * `key` - key.
    * `value` - value.
  * `is_dst_read_only` - destination readonly set.
  * `is_migrate_account` - migrate account.
  * `is_override_root` - override root destination by source database.
  * `migrate_type` - migrate type.
* `run_mode` - run mode.
* `src_info` - source info.
  * `access_type` - access type.
  * `database_type` - database type.
  * `extra_attr` - extra attributes.
    * `key` - key.
    * `value` - value.
  * `info` - databse info list.
    * `account_mode` - account mode.
    * `account_role` - account role.
    * `account` - account.
    * `ccn_gw_id` - ccn gateway id.
    * `cvm_instance_id` - cvm instance id.
    * `db_kernel` - database kernel.
    * `engine_version` - engine version.
    * `host` - host.
    * `instance_id` - instance id.
    * `password` - password.
    * `port` - port.
    * `role` - node role.
    * `subnet_id` - subnet id.
    * `tmp_secret_id` - temporary secret id.
    * `tmp_secret_key` - temporary secret key.
    * `tmp_token` - temporary token.
    * `uniq_vpn_gw_id` - vpn gateway id.
    * `user` - user.
    * `vpc_id` - vpc id.
  * `node_type` - node type.
  * `region` - region.
  * `supplier` - supplier.


## Import

dts migrate_job can be imported using the id, e.g.
```
$ terraform import tencentcloud_dts_migrate_job.migrate_job migrateJob_id
```

