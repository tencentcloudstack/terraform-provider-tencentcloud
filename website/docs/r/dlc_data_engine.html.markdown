---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_data_engine"
sidebar_current: "docs-tencentcloud-resource-dlc_data_engine"
description: |-
  Provides a resource to create a DLC data engine
---

# tencentcloud_dlc_data_engine

Provides a resource to create a DLC data engine

## Example Usage

```hcl
resource "tencentcloud_dlc_data_engine" "example" {
  engine_type        = "spark"
  data_engine_name   = "tf-example"
  cluster_type       = "spark_cu"
  mode               = 1
  auto_resume        = false
  size               = 16
  min_clusters       = 1
  max_clusters       = 1
  cidr_block         = "10.255.0.0/16"
  message            = "DLC data engine demo."
  image_version_name = "Standard-S 1.1"
  engine_exec_type   = "BATCH"
  engine_generation  = "Native"
  session_resource_template {
    driver_size          = "medium"
    executor_max_numbers = 7
    executor_nums        = 1
    executor_size        = "medium"
  }
}
```

## Argument Reference

The following arguments are supported:

* `auto_resume` - (Required, Bool) Whether to automatically start the clusters.
* `cluster_type` - (Required, String) The cluster type. Valid values: `spark_private`, `presto_private`, `presto_cu`, and `spark_cu`.
* `data_engine_name` - (Required, String) The name of the virtual cluster.
* `engine_type` - (Required, String) The engine type. Valid values: `spark` and `presto`.
* `mode` - (Required, Int) The billing mode. Valid values: `0` (shared engine), `1` (pay-as-you-go), and `2` (monthly subscription).
* `auto_authorization` - (Optional, Bool) Automatic authorization.
* `auto_renew` - (Optional, Int) The auto-renewal status of the resource. For the postpaid mode, no renewal is required, and the value is fixed to `0`. For the prepaid mode, valid values are `0` (manual), `1` (auto), and `2` (no renewal). If this parameter is set to `0` for a key account in the prepaid mode, auto-renewal applies. It defaults to `0`.
* `auto_suspend_time` - (Optional, Int) The cluster auto-suspension time, which defaults to 10 min.
* `auto_suspend` - (Optional, Bool) Whether to automatically suspend clusters. Valid values: `false` (default, no) and `true` (yes).
* `cidr_block` - (Optional, String) The VPC CIDR block.
* `crontab_resume_suspend_strategy` - (Optional, List) The complex policy for scheduled start and suspension, including the start/suspension time and suspension policy.
* `crontab_resume_suspend` - (Optional, Int) Whether to enable scheduled start and suspension of clusters. Valid values: `0` (disable) and `1` (enable). Note: This policy and the auto-suspension policy are mutually exclusive.
* `data_engine_config_pairs` - (Optional, List) The advanced configurations of clusters.
* `default_data_engine` - (Optional, Bool) Whether it is the default virtual cluster.
* `elastic_limit` - (Optional, Int) The upper limit (in CUs) for scaling of the monthly subscribed Spark job cluster.
* `elastic_switch` - (Optional, Bool) Whether to enable the scaling feature for a monthly subscribed Spark job cluster.
* `engine_exec_type` - (Optional, String) The type of tasks to be executed by the engine, which defaults to SQL. Valid values: `SQL` and `BATCH`.
* `engine_generation` - (Optional, String) Generation of the engine. SuperSQL means the supersql engine while Native means the standard engine. It is SuperSQL by default.
* `engine_network_id` - (Optional, String) Engine network ID.
* `image_version_name` - (Optional, String) The version name of cluster image, such as SuperSQL-P 1.1 and SuperSQL-S 3.2. If no value is passed in, a cluster is created using the latest image version.
* `main_cluster_name` - (Optional, String) The primary cluster, which is specified when a failover cluster is created.
* `max_clusters` - (Optional, Int) The maximum number of clusters.
* `max_concurrency` - (Optional, Int) The max task concurrency of a cluster, which defaults to 5.
* `message` - (Optional, String) The description.
* `min_clusters` - (Optional, Int) The minimum number of clusters.
* `pay_mode` - (Optional, Int) The pay mode. Valid value: `0` (postpaid, default) and `1` (prepaid) (currently not available).
* `resource_type` - (Optional, String) The resource type. Valid values: `Standard_CU` (standard) and `Memory_CU` (memory).
* `session_resource_template` - (Optional, List) The session resource configuration template for a Spark job cluster.
* `size` - (Optional, Int) Cluster size. Required when updating.
* `time_span` - (Optional, Int) The usage duration of the resource. Postpaid: Fill in 3,600 as a fixed figure; prepaid: fill in a figure equal to or bigger than 1 which means purchasing resources for one month. The maximum figure is not bigger than 120. The default value is 1.
* `time_unit` - (Optional, String) The unit of the resource period. Valid values: `s` (default) for the postpaid mode and `m` for the prepaid mode.
* `tolerable_queue_time` - (Optional, Int) The task queue time limit, which defaults to 0. When the actual queue time exceeds the value set here, scale-out may be triggered. Setting this parameter to 0 represents that scale-out may be triggered immediately after a task queues up.

The `crontab_resume_suspend_strategy` object supports the following:

* `resume_time` - (Optional, String) Scheduled starting time, such as 8: 00 a.m. on Monday and Wednesday.
* `suspend_strategy` - (Optional, Int) The suspension setting. Valid values: `0` (suspension after task end, default) and `1` (force suspension).
* `suspend_time` - (Optional, String) Scheduled suspension time, such as 8: 00 p.m. on Monday and Wednesday.

The `data_engine_config_pairs` object supports the following:

* `config_item` - (Required, String) Configuration items.
* `config_value` - (Required, String) Configuration value.

The `running_time_parameters` object of `session_resource_template` supports the following:

* `config_item` - (Required, String) Configuration items.
* `config_value` - (Required, String) Configuration value.

The `session_resource_template` object supports the following:

* `driver_size` - (Optional, String) The driver size. Valid values for the standard resource type: `small`, `medium`, `large`, and `xlarge`. Valid values for the memory resource type: `m.small`, `m.medium`, `m.large`, and `m.xlarge`.
* `executor_max_numbers` - (Optional, Int) The maximum executor count (in dynamic mode). The minimum value is 1 and the maximum value is less than the cluster specification. If you set `ExecutorMaxNumbers` to a value smaller than that of `ExecutorNums`, the value of `ExecutorMaxNumbers` is automatically changed to that of `ExecutorNums`.
* `executor_nums` - (Optional, Int) The executor count. The minimum value is 1 and the maximum value is less than the cluster specification.
* `executor_size` - (Optional, String) The executor size. Valid values for the standard resource type: `small`, `medium`, `large`, and `xlarge`. Valid values for the memory resource type: `m.small`, `m.medium`, `m.large`, and `m.xlarge`.
* `running_time_parameters` - (Optional, List) The running time parameters of the session resource configuration template for a Spark job cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `data_engine_id` - Data engine ID.


## Import

DLC data engine can be imported using the dataEngineName#dataEngineId, e.g.

```
terraform import tencentcloud_dlc_data_engine.example tf-example#DataEngine-d3gk8r5h
```

