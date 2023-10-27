---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_describe_data_engine"
sidebar_current: "docs-tencentcloud-datasource-dlc_describe_data_engine"
description: |-
  Use this data source to query detailed information of dlc describe_data_engine
---

# tencentcloud_dlc_describe_data_engine

Use this data source to query detailed information of dlc describe_data_engine

## Example Usage

```hcl
data "tencentcloud_dlc_describe_data_engine" "describe_data_engine" {
  data_engine_name = "testSpark"
}
```

## Argument Reference

The following arguments are supported:

* `data_engine_name` - (Required, String) Engine name.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data_engine` - Data engine details.
  * `auto_resume` - Whether to automatically start the cluster, prepay not support.
  * `auto_suspend_time` - Cluster automatic suspension time, default 10 minutes.
  * `auto_suspend` - Whether to automatically suspend the cluster, prepay not support.
  * `child_image_version_id` - Engine Image version id.
  * `cidr_block` - Engine VPC network segment, just like 192.0.2.1/24.
  * `cluster_type` - Engine cluster type, only support: spark_cu/presto_cu.
  * `create_time` - Create time.
  * `crontab_resume_suspend_strategy` - Engine auto suspend strategy, when AutoSuspend is true, CrontabResumeSuspend must stop.
    * `resume_time` - Scheduled pull-up time: For example: 8 o&amp;#39;clock on Monday is expressed as 1000000-08:00:00.
    * `suspend_strategy` - Suspend configuration: 0 (default): wait for the task to end before suspending, 1: force suspend.
    * `suspend_time` - Scheduled suspension time: For example: 20 o&amp;#39;clock on Monday is expressed as 1000000-20:00:00.
  * `crontab_resume_suspend` - Engine crontab resume or suspend strategy, only support: 0: Wait(default), 1: Kill.
  * `data_engine_id` - Engine unique id.
  * `data_engine_name` - Engine name.
  * `default_data_engine` - Whether it is the default virtual cluster.
  * `default_house` - Is it the default engine?.
  * `elastic_limit` - For spark Batch ExecType, yearly and monthly cluster elastic limit.
  * `elastic_switch` - For spark Batch ExecType, yearly and monthly cluster whether to enable elasticity.
  * `engine_exec_type` - Engine exec type, only support SQL(default) or BATCH.
  * `engine_type` - Engine type, only support: spark/presto.
  * `expire_time` - Expire time.
  * `image_version_id` - Engine major version id.
  * `image_version_name` - Engine image version name.
  * `isolated_time` - Isolated time.
  * `max_clusters` - Engine max cluster size,  MaxClusters less than or equal to 10 and MaxClusters bigger than MinClusters.
  * `max_concurrency` - Maximum number of concurrent tasks in a single cluster, default 5.
  * `message` - Engine description information.
  * `min_clusters` - Engine min size, greater than or equal to 1 and MaxClusters bigger than MinClusters.
  * `mode` - Engine mode, only support 1: ByAmount, 2: YearlyAndMonthly.
  * `network_connection_set` - Network connection configuration.
    * `appid` - User appid.
    * `associate_id` - Network configuration unique identifier.
    * `create_time` - Create time.
    * `datasource_connection_cidr_block` - Datasource connection cidr block.
    * `datasource_connection_id` - Data source id (obsolete).
    * `datasource_connection_name` - Network configuration name.
    * `datasource_connection_subnet_cidr_block` - Datasource connection subnet cidr block.
    * `datasource_connection_subnet_id` - Datasource subnetId.
    * `datasource_connection_vpc_id` - Datasource vpcid.
    * `house_id` - Data engine id.
    * `house_name` - Data engine name.
    * `id` - Network configuration id.
    * `network_connection_desc` - Network configuration description.
    * `network_connection_type` - Network configuration type.
    * `state` - Network configuration status (0-initialization, 1-normal).
    * `sub_account_uin` - User sub uin.
    * `uin` - User uin.
    * `update_time` - Update time.
  * `permissions` - Engine permissions.
  * `quota_id` - Reference ID.
  * `renew_flag` - Automatic renewal flag, 0, initial state, automatic renewal is not performed by default. If the user has prepaid non-stop service privileges, automatic renewal will occur. 1: Automatic renewal. 2: Make it clear that there will be no automatic renewal.
  * `resource_type` - Engine resource type not match, only support: Standard_CU/Memory_CU(only BATCH ExecType).
  * `reversal_time` - Reversal time.
  * `session_resource_template` - For spark Batch ExecType, cluster session resource configuration template.
    * `driver_size` - Engine driver size specification only supports: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge.
    * `executor_max_numbers` - Specify the executor max number (in a dynamic configuration scenario), the minimum value is 1, and the maximum value is less than the cluster specification (when ExecutorMaxNumbers is less than ExecutorNums, the value is set to ExecutorNums).
    * `executor_nums` - Specify the number of executors. The minimum value is 1 and the maximum value is less than the cluster specification.
    * `executor_size` - Engine executor size specification only supports: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge.
  * `size` - Engine size.
  * `spend_after` - Automatic recovery time, prepay not support.
  * `start_standby_cluster` - Whether to enable the backup cluster.
  * `state` - Engine state, only support: 0:Init/-1:Failed/-2:Deleted/1:Pause/2:Running/3:ToBeDelete/4:Deleting.
  * `sub_account_uin` - Operator.
  * `tag_list` - Tag list.
    * `tag_key` - Tag key.
    * `tag_value` - Tag value.
  * `tolerable_queue_time` - Tolerable queuing time, default 0. scaling may be triggered when tasks are queued for longer than the tolerable time. if this parameter is 0, it means that capacity expansion may be triggered immediately once a task is queued.
  * `ui_u_r_l` - Jump address of ui.
  * `update_time` - Update time.
  * `user_alias` - User&amp;#39;s nickname.
  * `user_app_id` - User appid.
  * `user_uin` - User uin.


