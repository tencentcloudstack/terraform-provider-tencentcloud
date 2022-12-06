---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_migrate_jobs"
sidebar_current: "docs-tencentcloud-datasource-dts_migrate_jobs"
description: |-
  Use this data source to query detailed information of dts migrateJobs
---

# tencentcloud_dts_migrate_jobs

Use this data source to query detailed information of dts migrateJobs

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

data "tencentcloud_dts_migrate_jobs" "all" {}

data "tencentcloud_dts_migrate_jobs" "job" {
  job_id   = tencentcloud_dts_migrate_job.migrate_job.id
  job_name = tencentcloud_dts_migrate_job.migrate_job.job_name
  status   = ["created"]
}

data "tencentcloud_dts_migrate_jobs" "src_dest" {

  src_region        = "ap-guangzhou"
  src_database_type = ["mysql"]
  dst_region        = "ap-guangzhou"
  dst_database_type = ["cynosdbmysql"]

  status = ["created"]
  tag_filters {
    tag_key   = "aaa"
    tag_value = "bbb"
  }
}
```

## Argument Reference

The following arguments are supported:

* `dst_access_type` - (Optional, Set: [`String`]) destination access type.
* `dst_database_type` - (Optional, Set: [`String`]) destination database type.
* `dst_instance_id` - (Optional, String) source instance id.
* `dst_region` - (Optional, String) destination region.
* `job_id` - (Optional, String) job id.
* `job_name` - (Optional, String) job name.
* `order_seq` - (Optional, String) order by, default by create time.
* `result_output_file` - (Optional, String) Used to save results.
* `run_mode` - (Optional, String) run mode.
* `src_access_type` - (Optional, Set: [`String`]) source access type.
* `src_database_type` - (Optional, Set: [`String`]) source database type.
* `src_instance_id` - (Optional, String) source instance id.
* `src_region` - (Optional, String) source region.
* `status` - (Optional, Set: [`String`]) migrate status.
* `tag_filters` - (Optional, List) tag filters.

The `tag_filters` object supports the following:

* `tag_key` - (Optional, String) tag key.
* `tag_value` - (Optional, String) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - migration job list.
  * `action` - action info.
    * `all_action` - all action list.
    * `allowed_action` - allowed action list.
  * `brief_msg` - brief message for migrate error.
  * `compare_task` - compare task info.
    * `compare_task_id` - compare task id.
    * `status` - status.
  * `create_time` - create time.
  * `dst_info` - destination info.
    * `access_type` - access type.
    * `database_type` - database type.
    * `info` - db info.
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
  * `end_time` - end time.
  * `expect_run_time` - expected run time.
  * `job_id` - job id.
  * `job_name` - job name.
  * `src_info` - source info.
    * `access_type` - access type.
    * `database_type` - database type.
    * `extra_attr` - extra attributes.
      * `key` - key.
      * `value` - value.
    * `info` - db info.
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
  * `start_time` - start time.
  * `status` - status.
  * `step_info` - step info.
    * `master_slave_distance` - master slave distance.
    * `seconds_behind_master` - seconds behind master.
    * `step_all` - number of all steps.
    * `step_info` - step infos.
      * `errors` - error list.
        * `help_doc` - help document.
        * `message` - message.
        * `solution` - solution.
      * `percent` - the percent of miragtion progress.
      * `start_time` - start time.
      * `status` - current status.
      * `step_id` - step id.
      * `step_message` - step message.
      * `step_name` - step name.
      * `step_no` - step number.
      * `warnings` - warning list.
        * `help_doc` - help document.
        * `message` - message.
        * `solution` - solution.
    * `step_now` - current step.
  * `tags` - tag list.
    * `tag_key` - tag key.
    * `tag_value` - tag value.
  * `trade_info` - trade info.
    * `billing_type` - billing type.
    * `deal_name` - deal name.
    * `expire_time` - expired time.
    * `instance_class` - instance class.
    * `isolate_reason` - isolate reason.
    * `isolate_time` - isolate time.
    * `last_deal_name` - last deal name.
    * `offline_reason` - offline reason.
    * `offline_time` - offline time.
    * `pay_type` - pay type.
    * `trade_status` - trade status.
  * `update_time` - update time.


