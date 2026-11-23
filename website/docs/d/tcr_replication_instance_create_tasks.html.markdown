---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_replication_instance_create_tasks"
sidebar_current: "docs-tencentcloud-datasource-tcr_replication_instance_create_tasks"
description: |-
  Use this data source to query detailed information of tcr replication instance create tasks
---

# tencentcloud_tcr_replication_instance_create_tasks

Use this data source to query detailed information of tcr replication instance create tasks

## Example Usage

```hcl
data "tencentcloud_tcr_replication_instance_create_tasks" "example" {
  replication_registry_id = "tcr-dg284imq-5-os5gcu"
  replication_region_id   = 5
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

## API Reference

This data source queries the following TencentCloud APIs:

* Read: [DescribeReplicationInstanceCreateTasks](https://www.tencentcloud.com/document/api/1051/42672)

