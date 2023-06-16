---
subcategory: "Tencent Container Registry(TCR)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_tcr_replication_instance_sync_status"
sidebar_current: "docs-tencentcloud-datasource-tcr_replication_instance_sync_status"
description: |-
  Use this data source to query detailed information of tcr replication_instance_sync_status
---

# tencentcloud_tcr_replication_instance_sync_status

Use this data source to query detailed information of tcr replication_instance_sync_status

## Example Usage

```hcl
data "tencentcloud_tcr_replication_instance_sync_status" "sync_status" {
  registry_id             = local.src_registry_id
  replication_registry_id = local.dst_registry_id
  replication_region_id   = local.dst_region_id
  show_replication_log    = false
}
```

## Argument Reference

The following arguments are supported:

* `registry_id` - (Required, String) master registry id.
* `replication_registry_id` - (Required, String) synchronization instance id.
* `replication_region_id` - (Optional, Int) synchronization instance region id.
* `result_output_file` - (Optional, String) Used to save results.
* `show_replication_log` - (Optional, Bool) whether to display the synchronization log.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `replication_log` - sync log. Note: This field may return null, indicating that no valid value can be obtained.
  * `destination` - destination resource. Note: This field may return null, indicating that no valid value can be obtained.
  * `end_time` - end time. Note: This field may return null, indicating that no valid value can be obtained.
  * `resource_type` - resource type. Note: This field may return null, indicating that no valid value can be obtained.
  * `source` - Source image. Note: This field may return null, indicating that no valid value can be obtained.
  * `start_time` - start time. Note: This field may return null, indicating that no valid value can be obtained.
  * `status` - sync status. Note: This field may return null, indicating that no valid value can be obtained.
* `replication_status` - sync status.
* `replication_time` - sync complete time.


