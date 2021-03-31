---
subcategory: "CVM Dedicated Host(CDH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cdh_instances"
sidebar_current: "docs-tencentcloud-datasource-cdh_instances"
description: |-
  Use this data source to query CDH instances.
---

# tencentcloud_cdh_instances

Use this data source to query CDH instances.

## Example Usage

```hcl
data "tencentcloud_cdh_instances" "list" {
  availability_zone = "ap-guangzhou-3"
  host_id           = "host-d6s7i5q4"
  host_name         = "test"
  host_state        = "RUNNING"
  project_id        = 1154137
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) The available zone that the CDH instance locates at.
* `host_id` - (Optional) ID of the CDH instances to be queried.
* `host_name` - (Optional) Name of the CDH instances to be queried.
* `host_state` - (Optional) State of the CDH instances to be queried. Valid values: `PENDING`, `LAUNCH_FAILURE`, `RUNNING`, `EXPIRED`.
* `project_id` - (Optional) The project CDH belongs to.
* `result_output_file` - (Optional) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `cdh_instance_list` - An information list of cdh instance. Each element contains the following attributes:
  * `availability_zone` - The available zone that the CDH instance locates at.
  * `cage_id` - Cage ID of the CDH instance. This parameter is only valid for CDH instances in the cages of finance availability zones.
  * `charge_type` - The charge type of the CDH instance.
  * `create_time` - Creation time of the CDH instance.
  * `cvm_instance_ids` - Id of CVM instances that have been created on the CDH instance.
  * `expired_time` - Expired time of the CDH instance.
  * `host_id` - ID of the CDH instance.
  * `host_name` - Name of the CDH instance.
  * `host_resource` - An information list of host resource. Each element contains the following attributes:
    * `cpu_available_num` - The number of available CPU cores of the instance.
    * `cpu_total_num` - The number of total CPU cores of the instance.
    * `disk_available_size` - Instance disk available capacity, unit in GiB.
    * `disk_total_size` - Instance disk total capacity, unit in GiB.
    * `disk_type` - Type of the disk.
    * `memory_available_size` - Instance memory available capacity, unit in GiB.
    * `memory_total_size` - Instance memory total capacity, unit in GiB.
  * `host_state` - State of the CDH instance.
  * `host_type` - Type of the CDH instance.
  * `prepaid_renew_flag` - Auto renewal flag.
  * `project_id` - The project CDH belongs to.


