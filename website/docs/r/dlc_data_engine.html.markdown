---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_data_engine"
sidebar_current: "docs-tencentcloud-resource-dlc_data_engine"
description: |-
  Provides a resource to create a dlc data_engine
---

# tencentcloud_dlc_data_engine

Provides a resource to create a dlc data_engine

## Example Usage

```hcl
resource "tencentcloud_dlc_data_engine" "data_engine" {
  engine_type            = "spark"
  data_engine_name       = "testSpark"
  cluster_type           = "spark_cu"
  mode                   = 1
  auto_resume            = false
  size                   = 16
  pay_mode               = 0
  min_clusters           = 1
  max_clusters           = 1
  default_data_engine    = false
  cidr_block             = "10.255.0.0/16"
  message                = "test spark1"
  time_span              = 1
  time_unit              = "h"
  auto_suspend           = false
  crontab_resume_suspend = 0
  engine_exec_type       = "BATCH"
}
```

## Argument Reference

The following arguments are supported:

* `auto_resume` - (Required, Bool) Whether to automatically start the cluster, prepay not support.
* `cluster_type` - (Required, String) Engine cluster type, only support: spark_cu/presto_cu.
* `data_engine_name` - (Required, String) Engine name.
* `engine_type` - (Required, String) Engine type, only support: spark/presto.
* `mode` - (Required, Int) Engine mode, only support 1: ByAmount, 2: YearlyAndMonthly.
* `auto_renew` - (Optional, Int) Engine auto renew, only support 0: Default, 1: AutoRenewON, 2: AutoRenewOFF.
* `auto_suspend_time` - (Optional, Int) Cluster automatic suspension time, default 10 minutes.
* `auto_suspend` - (Optional, Bool) Whether to automatically suspend the cluster, prepay not support.
* `cidr_block` - (Optional, String) Engine VPC network segment, just like 192.0.2.1/24.
* `crontab_resume_suspend_strategy` - (Optional, List) Engine auto suspend strategy, when AutoSuspend is true, CrontabResumeSuspend must stop.
* `crontab_resume_suspend` - (Optional, Int) Engine crontab resume or suspend strategy, only support: 0: Wait(default), 1: Kill.
* `data_engine_config_pairs` - (Optional, List) Cluster advanced configuration.
* `default_data_engine` - (Optional, Bool) Whether it is the default virtual cluster.
* `elastic_limit` - (Optional, Int) For spark Batch ExecType, yearly and monthly cluster elastic limit.
* `elastic_switch` - (Optional, Bool) For spark Batch ExecType, yearly and monthly cluster whether to enable elasticity.
* `engine_exec_type` - (Optional, String) Engine exec type, only support SQL(default) or BATCH.
* `image_version_name` - (Optional, String) Cluster image version name. Such as SuperSQL-P 1.1; SuperSQL-S 3.2, etc., do not upload, and create a cluster with the latest mirror version by default.
* `main_cluster_name` - (Optional, String) Primary cluster name, specified when creating a disaster recovery cluster.
* `max_clusters` - (Optional, Int) Engine max cluster size, MaxClusters less than or equal to 10 and MaxClusters bigger than MinClusters.
* `max_concurrency` - (Optional, Int) Maximum number of concurrent tasks in a single cluster, default 5.
* `message` - (Optional, String) Engine description information.
* `min_clusters` - (Optional, Int) Engine min size, greater than or equal to 1 and MaxClusters bigger than MinClusters.
* `pay_mode` - (Optional, Int) Engine pay mode type, only support 0: postPay, 1: prePay(default).
* `resource_type` - (Optional, String) Engine resource type not match, only support: Standard_CU/Memory_CU(only BATCH ExecType).
* `session_resource_template` - (Optional, List) For spark Batch ExecType, cluster session resource configuration template.
* `size` - (Optional, Int) Cluster size. Required when updating.
* `time_span` - (Optional, Int) Engine TimeSpan, prePay: minimum of 1, representing one month of purchasing resources, with a maximum of 120, default 3600, postPay: fixed fee of 3600.
* `time_unit` - (Optional, String) Engine TimeUnit, prePay: use m(default), postPay: use h.
* `tolerable_queue_time` - (Optional, Int) Tolerable queuing time, default 0. scaling may be triggered when tasks are queued for longer than the tolerable time. if this parameter is 0, it means that capacity expansion may be triggered immediately once a task is queued.

The `crontab_resume_suspend_strategy` object supports the following:

* `resume_time` - (Optional, String) Scheduled pull-up time: For example: 8 o&amp;#39;clock on Monday is expressed as 1000000-08:00:00.
* `suspend_strategy` - (Optional, Int) Suspend configuration: 0 (default): wait for the task to end before suspending, 1: force suspend.
* `suspend_time` - (Optional, String) Scheduled suspension time: For example: 20 o&amp;#39;clock on Monday is expressed as 1000000-20:00:00.

The `data_engine_config_pairs` object supports the following:

* `config_item` - (Required, String) Configuration items.
* `config_value` - (Required, String) Configuration value.

The `session_resource_template` object supports the following:

* `driver_size` - (Optional, String) Engine driver size specification only supports: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge.
* `executor_max_numbers` - (Optional, Int) Specify the executor max number (in a dynamic configuration scenario), the minimum value is 1, and the maximum value is less than the cluster specification (when ExecutorMaxNumbers is less than ExecutorNums, the value is set to ExecutorNums).
* `executor_nums` - (Optional, Int) Specify the number of executors. The minimum value is 1 and the maximum value is less than the cluster specification.
* `executor_size` - (Optional, String) Engine executor size specification only supports: small/medium/large/xlarge/m.small/m.medium/m.large/m.xlarge.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dlc data_engine can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_data_engine.data_engine data_engine_id
```

