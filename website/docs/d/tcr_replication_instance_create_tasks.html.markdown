---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_replication_instance_create_tasks"
sidebar_current: "docs-tencentcloud-datasource-tcr_replication_instance_create_tasks"
description: |-
  Use this data source to query detailed information of tcr replication_instance_create_tasks
---

# tencentcloud_tcr_replication_instance_create_tasks

Use this data source to query detailed information of tcr replication_instance_create_tasks

## Example Usage

```hcl
local {
  src_registry_id = local.tcr_id
  dst_registry_id = tencentcloud_tcr_manage_replication_operation.my_replica.destination_registry_id
  dst_region_id   = tencentcloud_tcr_manage_replication_operation.my_replica.destination_region_id
}

data "tencentcloud_tcr_replication_instance_create_tasks" "create_tasks" {
  replication_registry_id = local.dst_registry_id
  replication_region_id   = local.dst_region_id
}
```

## Argument Reference

The following arguments are supported:

* `replication_region_id` - (Required, Int) synchronization instance region Id, see ReplicationRegionId in DescribeReplicationInstances.
* `replication_registry_id` - (Required, String) synchronization instance Id, see RegistryId in DescribeReplicationInstances.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `status` - overall task status.
* `task_detail` - task details.
  * `created_time` - task start name.
  * `finished_time` - task end time. Note: This field may return null, indicating that no valid value can be obtained.
  * `task_message` - Task status information. Note: This field may return null, indicating that no valid value can be obtained.
  * `task_name` - task name.
  * `task_status` - task status.
  * `task_uuid` - task UUID.


