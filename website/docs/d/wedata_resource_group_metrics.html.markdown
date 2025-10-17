---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_resource_group_metrics"
sidebar_current: "docs-tencentcloud-datasource-wedata_resource_group_metrics"
description: |-
  Use this data source to query detailed information of WeData resource group metrics
---

# tencentcloud_wedata_resource_group_metrics

Use this data source to query detailed information of WeData resource group metrics

## Example Usage

```hcl
data "tencentcloud_wedata_resource_group_metrics" "example" {
  resource_group_id = "20250909193110713075"
}
```

## Argument Reference

The following arguments are supported:

* `resource_group_id` - (Required, String) Execution resource group ID.
* `end_time` - (Optional, Int) Usage trend end time (milliseconds), default to current time.
* `granularity` - (Optional, Int) Metric collection granularity, unit in minutes, default 1 minute.
* `metric_type` - (Optional, String) Metric dimension.

- all --- All
- task --- Task metrics
- system --- System metrics.
* `result_output_file` - (Optional, String) Used to save results.
* `start_time` - (Optional, Int) Usage trend start time (milliseconds), default to the last hour.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - Execution group metric information.
  * `cpu_num` - Resource group specification related: CPU count.
  * `disk_volume` - Resource group specification related: disk specification.
  * `life_cycle` - Resource group lifecycle, unit: days.
  * `maximum_concurrency` - Resource group specification related: maximum concurrency.
  * `mem_size` - Resource group specification related: memory size, unit: G.
  * `metric_snapshots` - Metric details.
    * `metric_name` - Metric name.

- ConcurrencyUsage --- Concurrency usage rate
- CpuCoreUsage --- CPU usage rate
- CpuLoad --- CPU load
- DevelopQueueTask --- Number of development tasks in queue
- DevelopRunningTask --- Number of running development tasks
- DevelopSchedulingTask --- Number of scheduling development tasks
- DiskUsage --- Disk usage
- DiskUsed --- Disk used amount
- MaximumConcurrency --- Maximum concurrency
- MemoryLoad --- Memory load
- MemoryUsage --- Memory usage.
    * `snapshot_value` - Current value.
    * `trend_list` - Metric trend.
      * `timestamp` - Timestamp.
      * `value` - Metric value.
  * `status` - Resource group status.

- 0 --- Initializing
- 1 --- Running
- 2 --- Running abnormally
- 3 --- Releasing
- 4 --- Released
- 5 --- Creating
- 6 --- Creation failed
- 7 --- Updating
- 8 --- Update failed
- 9 --- Expired
- 10 --- Release failed
- 11 --- In use
- 12 --- Not in use.


