---
subcategory: "TDSQL-C MySQL(CynosDB)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cynosdb_readonly_instance"
sidebar_current: "docs-tencentcloud-resource-cynosdb_readonly_instance"
description: |-
  Provide a resource to create a CynosDB readonly instance.
---

# tencentcloud_cynosdb_readonly_instance

Provide a resource to create a CynosDB readonly instance.

## Example Usage

```hcl
resource "tencentcloud_cynosdb_readonly_instance" "foo" {
  cluster_id           = cynosdbmysql-dzj5l8gz
  instance_name        = "tf-cynosdb-readonly-instance"
  force_delete         = true
  instance_cpu_core    = 2
  instance_memory_size = 4

  instance_maintain_duration   = 7200
  instance_maintain_start_time = 21600
  instance_maintain_weekdays = [
    "Fri",
    "Mon",
    "Sat",
    "Sun",
    "Thu",
    "Wed",
    "Tue",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required, String, ForceNew) Cluster ID which the readonly instance belongs to.
* `instance_name` - (Required, String) Instance name.
* `force_delete` - (Optional, Bool) Indicate whether to delete readonly instance directly or not. Default is false. If set true, instance will be deleted instead of staying recycle bin. Note: works for both `PREPAID` and `POSTPAID_BY_HOUR` cluster.
* `instance_cpu_core` - (Optional, Int) The number of CPU cores of read-write type instance in the CynosDB cluster. Required while creating normal cluster. Note: modification of this field will take effect immediately, if want to upgrade on maintenance window, please upgrade from console.
* `instance_maintain_duration` - (Optional, Int) Duration time for maintenance, unit in second. `3600` by default.
* `instance_maintain_start_time` - (Optional, Int) Offset time from 00:00, unit in second. For example, 03:00am should be `10800`. `10800` by default.
* `instance_maintain_weekdays` - (Optional, Set: [`String`]) Weekdays for maintenance. `["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"]` by default.
* `instance_memory_size` - (Optional, Int) Memory capacity of read-write type instance, unit in GB. Required while creating normal cluster. Note: modification of this field will take effect immediately, if want to upgrade on maintenance window, please upgrade from console.
* `subnet_id` - (Optional, String) ID of the subnet within this VPC.
* `vpc_id` - (Optional, String) ID of the VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `instance_status` - Status of the instance.
* `instance_storage_size` - Storage size of the instance, unit in GB.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [AddInstances](https://cloud.tencent.com/document/api/1003/48088) - This API is used to add instances to a cluster.
* Create: [DescribeResourcesByDealName](https://cloud.tencent.com/document/api/1003/52131) - This interface (DescribeResourcesByDealName) is used to query order-associated instances.
* Create: [ModifyMaintainPeriodConfig](https://cloud.tencent.com/document/api/1003/48080) - This API is used to modify maintenance time configuration.
* Read: [DescribeInstanceDetail](https://cloud.tencent.com/document/api/1003/48250) - This API is used to query instance details.
* Read: [DescribeInstances](https://cloud.tencent.com/document/api/1003/48334) - This API is used to query the list of instances.
* Read: [DescribeMaintainPeriod](https://cloud.tencent.com/document/api/1003/48083) - This interface (DescribeMaintainPeriod) is used to query the instance maintenance window.
* Update: [DescribeInstanceDetail](https://cloud.tencent.com/document/api/1003/48250) - This API is used to query instance details.
* Update: [DescribeInstances](https://cloud.tencent.com/document/api/1003/48334) - This API is used to query the list of instances.
* Update: [ModifyInstanceName](https://cloud.tencent.com/document/api/1003/71728) - This API is used to modify instance name.
* Update: [ModifyMaintainPeriodConfig](https://cloud.tencent.com/document/api/1003/48080) - This API is used to modify maintenance time configuration.
* Update: [UpgradeInstance](https://cloud.tencent.com/document/api/1003/48077) - This interface (UpgradeInstance) is used to upgrade instances.
* Delete: [DescribeInstanceDetail](https://cloud.tencent.com/document/api/1003/48250) - This API is used to query instance details.
* Delete: [DescribeInstances](https://cloud.tencent.com/document/api/1003/48334) - This API is used to query the list of instances.
* Delete: [IsolateInstance](https://cloud.tencent.com/document/api/1003/48081) - This API is used to isolate an instance.
* Delete: [OfflineInstance](https://cloud.tencent.com/document/api/1003/48249) - This interface (OfflineInstance) is used to terminate an instance.

## Import

CynosDB readonly instance can be imported using the id, e.g.

```
$ terraform import tencentcloud_cynosdb_readonly_instance.foo cynosdbmysql-ins-dhwynib6
```

