---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_repair_tasks"
sidebar_current: "docs-tencentcloud-datasource-cvm_repair_tasks"
description: |-
  Use this data source to query CVM repair tasks.
---

# tencentcloud_cvm_repair_tasks

Use this data source to query CVM repair tasks.

## Example Usage

```hcl
data "tencentcloud_cvm_repair_tasks" "tasks" {
  task_status = [1, 4]
}
```

### Query with multiple filters

```hcl
data "tencentcloud_cvm_repair_tasks" "filtered" {
  product      = "CVM"
  task_status  = [1, 2]
  instance_ids = ["ins-xxxxxxxx"]
  start_date   = "2023-01-01 00:00:00"
  end_date     = "2023-12-31 23:59:59"
  order_field  = "CreateTime"
  order        = 1
}
```

## Argument Reference

The following arguments are supported:

* `aliases` - (Optional, Set: [`String`]) Instance name list (alias). Query tasks by instance names.
* `end_date` - (Optional, String) Query end date, format: YYYY-MM-DD hh:mm:ss. Filter tasks created until this date.
* `instance_ids` - (Optional, Set: [`String`]) Instance ID list (e.g., ins-xxxxxxxx). Query tasks by instance IDs.
* `order_field` - (Optional, String) Sorting field. Valid values: CreateTime (creation time), AuthTime (authorization time), EndTime (end time).
* `order` - (Optional, Int) Sorting order. 0: ascending, 1: descending. Default: 0.
* `product` - (Optional, String) Product type, optional values: CVM (Cloud Virtual Machine), CDH (Cloud Dedicated Host), CPM2.0 (Cloud Physical Machine 2.0).
* `result_output_file` - (Optional, String) Used to save results.
* `start_date` - (Optional, String) Query start date, format: YYYY-MM-DD hh:mm:ss. Filter tasks created from this date.
* `task_ids` - (Optional, Set: [`String`]) Task ID list (e.g., rep-xxxxxxxx). Query specific tasks by task IDs.
* `task_status` - (Optional, Set: [`Int`]) Task status list. Valid values: 1 (pending authorization), 2 (processing), 3 (ended), 4 (scheduled), 5 (cancelled), 6 (avoided).
* `task_type_ids` - (Optional, Set: [`Int`]) Task type ID list. Valid values: 101 (instance running hazard), 102 (instance running exception), 103 (instance hard disk exception), 104 (instance network connection exception), 105 (instance running warning), 106 (instance hard disk warning), 107 (instance maintenance upgrade).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `repair_task_list` - An information list of repair tasks. Each element contains the following attributes:
  * `alias` - Instance name (alias).
  * `auth_source` - Authorization source.
  * `auth_time` - Task authorization time.
  * `auth_type` - Authorization type.
  * `create_time` - Task creation time.
  * `device_status` - Device status.
  * `end_time` - Task end time.
  * `instance_id` - Instance ID.
  * `lan_ip` - Private IP address.
  * `operate_status` - Operation status.
  * `product` - Product type.
  * `region` - Region.
  * `subnet_id` - Subnet ID.
  * `subnet_name` - Subnet name.
  * `task_detail` - Task detail description.
  * `task_id` - Task ID.
  * `task_status` - Task status.
  * `task_sub_type` - Task sub-type.
  * `task_type_id` - Task type ID.
  * `task_type_name` - Task type name.
  * `vpc_id` - VPC ID.
  * `vpc_name` - VPC name.
  * `wan_ip` - Public IP address.
  * `zone` - Availability zone.
* `total_count` - Total count of repair tasks that match the filter conditions.


