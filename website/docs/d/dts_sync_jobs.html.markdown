---
subcategory: "Data Transmission Service(DTS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dts_sync_jobs"
sidebar_current: "docs-tencentcloud-datasource-dts_sync_jobs"
description: |-
  Use this data source to query detailed information of dts syncJobs
---

# tencentcloud_dts_sync_jobs

Use this data source to query detailed information of dts syncJobs

## Example Usage

```hcl
resource "tencentcloud_dts_sync_job" "job" {
  job_name          = "tf_dts_test"
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

data "tencentcloud_dts_sync_jobs" "sync_jobs" {
  job_id   = tencentcloud_dts_sync_job.job.id
  job_name = "tf_dts_test"
}
```

## Argument Reference

The following arguments are supported:

* `job_id` - (Optional, String) job id.
* `job_name` - (Optional, String) job name.
* `job_type` - (Optional, String) job type.
* `order_seq` - (Optional, String) order way, optional value is DESC or ASC.
* `order` - (Optional, String) order field.
* `pay_mode` - (Optional, String) pay mode, optional value is PrePay or PostPay.
* `result_output_file` - (Optional, String) Used to save results.
* `run_mode` - (Optional, String) run mode, optional value is mmediate or Timed.
* `status` - (Optional, Set: [`String`]) status.
* `tag_filters` - (Optional, List) tag filters.

The `tag_filters` object supports the following:

* `tag_key` - (Optional, String) tag key.
* `tag_value` - (Optional, String) tag value.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `list` - sync job list.
  * `actions` - support action list for current status.
  * `all_actions` - all action list.
  * `create_time` - create time.
  * `detail` - tag list.
    * `current_step_progress` - current step progress.
    * `master_slave_distance` - master slave distance.
    * `message` - message.
    * `progress` - progress.
    * `seconds_behind_master` - seconds behind master.
    * `step_all` - total step numbers.
    * `step_infos` - step infos.
      * `errors` - error list.
        * `code` - code.
        * `help_doc` - help document.
        * `message` - message.
        * `solution` - solution.
      * `progress` - current step progress.
      * `start_time` - start time.
      * `status` - current status.
      * `step_id` - step id.
      * `step_name` - step name.
      * `step_no` - step number.
      * `warnings` - waring list.
        * `code` - code.
        * `help_doc` - help document.
        * `message` - message.
        * `solution` - solution.
    * `step_now` - current step number.
  * `dst_database_type` - destination database type.
  * `dst_info` - destination info.
    * `account_mode` - account mode.
    * `account_role` - account role.
    * `account` - account.
    * `ccn_id` - ccn id.
    * `cvm_instance_id` - cvm instance id.
    * `db_kernel` - database kernel.
    * `db_name` - database name.
    * `engine_version` - engine version.
    * `instance_id` - instance id.
    * `ip` - ip.
    * `password` - password.
    * `port` - port.
    * `region` - region.
    * `subnet_id` - subnet id.
    * `supplier` - supplier.
    * `tmp_secret_id` - temporary secret id.
    * `tmp_secret_key` - temporary secret key.
    * `tmp_token` - temporary token.
    * `uniq_dcg_id` - dedicated gateway id.
    * `uniq_vpn_gw_id` - vpn gateway id.
    * `user` - user.
    * `vpc_id` - vpc id.
  * `dst_region` - destination region.
  * `end_time` - end time.
  * `expect_run_time` - expected run time.
  * `expire_time` - expire time.
  * `job_id` - job id.
  * `job_name` - job name.
  * `objects` - objects.
    * `advanced_objects` - advanced objects.
    * `databases` - database list.
      * `db_mode` - database mode.
      * `db_name` - database name.
      * `function_mode` - function mode.
      * `functions` - functions.
      * `new_db_name` - new database name.
      * `new_schema_name` - new schema name.
      * `procedure_mode` - procedure mode.
      * `procedures` - procedures.
      * `schema_name` - schema name.
      * `table_mode` - table mode.
      * `tables` - table list.
        * `filter_condition` - filter condition.
        * `new_table_name` - new table name.
        * `table_name` - table name.
      * `view_mode` - view mode.
      * `views` - view list.
        * `new_view_name` - new view name.
        * `view_name` - view name.
    * `mode` - object mode.
  * `options` - options.
    * `add_additional_column` - add additional column.
    * `conflict_handle_option` - conflict handle option.
      * `condition_column` - condition column.
      * `condition_operator` - condition override operator.
      * `condition_order_in_src_and_dst` - condition order in source and destination.
    * `conflict_handle_type` - conflict handle type.
    * `ddl_options` - ddl options.
      * `ddl_object` - ddl object.
      * `ddl_value` - ddl value.
    * `deal_of_exist_same_table` - deal of exist same table.
    * `init_type` - init type.
    * `op_types` - operation types.
  * `pay_mode` - pay mode.
  * `run_mode` - run mode.
  * `specification` - specification.
  * `src_access_type` - source access type.
  * `src_database_type` - source database type.
  * `src_info` - source info.
    * `account_mode` - account mode.
    * `account_role` - account role.
    * `account` - account.
    * `ccn_id` - ccn id.
    * `cvm_instance_id` - cvm instance id.
    * `db_kernel` - database kernel.
    * `db_name` - database name.
    * `engine_version` - engine version.
    * `instance_id` - instance id.
    * `ip` - ip.
    * `password` - password.
    * `port` - port.
    * `region` - region.
    * `subnet_id` - subnet id.
    * `supplier` - supplier.
    * `tmp_secret_id` - temporary secret id.
    * `tmp_secret_key` - temporary secret key.
    * `tmp_token` - temporary token.
    * `uniq_dcg_id` - dedicated gateway id.
    * `uniq_vpn_gw_id` - vpn gateway id.
    * `user` - user.
    * `vpc_id` - vpc id.
  * `src_region` - source region.
  * `start_time` - start time.
  * `status` - status.
  * `tags` - tag list.
    * `tag_key` - tag key.
    * `tag_value` - tag value.


