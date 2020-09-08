---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_snapshots"
sidebar_current: "docs-tencentcloud-datasource-cbs_snapshots"
description: |-
  Use this data source to query detailed information of CBS snapshots.
---

# tencentcloud_cbs_snapshots

Use this data source to query detailed information of CBS snapshots.

## Example Usage

```hcl
data "tencentcloud_cbs_snapshots" "snapshots" {
  snapshot_id        = "snap-f3io7adt"
  result_output_file = "mytestpath"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) The available zone that the CBS instance locates at.
* `project_id` - (Optional) ID of the project within the snapshot.
* `result_output_file` - (Optional) Used to save results.
* `snapshot_id` - (Optional) ID of the snapshot to be queried.
* `snapshot_name` - (Optional) Name of the snapshot to be queried.
* `storage_id` - (Optional) ID of the the CBS which this snapshot created from.
* `storage_usage` - (Optional) Types of CBS which this snapshot created from, and available values include SYSTEM_DISK and DATA_DISK.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `snapshot_list` - A list of snapshot. Each element contains the following attributes:
  * `availability_zone` - The available zone that the CBS instance locates at.
  * `create_time` - Creation time of snapshot.
  * `encrypt` - Indicates whether the snapshot is encrypted.
  * `percent` - Snapshot creation progress percentage.
  * `project_id` - ID of the project within the snapshot.
  * `snapshot_id` - ID of the snapshot.
  * `snapshot_name` - Name of the snapshot.
  * `storage_id` - ID of the the CBS which this snapshot created from.
  * `storage_size` - Volume of storage which this snapshot created from.
  * `storage_usage` - Types of CBS which this snapshot created from.


