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

### Query all snapshots

```hcl
data "tencentcloud_cbs_snapshots" "snapshots" {}
```

### Query snapshots by filters

```hcl
data "tencentcloud_cbs_snapshots" "snapshots" {
  snapshot_id        = "snap-hibh08s3"
  result_output_file = "my_snapshots"
}

data "tencentcloud_cbs_snapshots" "snapshots" {
  snapshot_name = "tf-example"
}

data "tencentcloud_cbs_snapshots" "snapshots" {
  storage_id = "disk-12j0fk1w"
}

data "tencentcloud_cbs_snapshots" "snapshots" {
  storage_usage = "SYSTEM_DISK"
}

data "tencentcloud_cbs_snapshots" "snapshots" {
  project_id = "0"
}

data "tencentcloud_cbs_snapshots" "snapshots" {
  availability_zone = "ap-guangzhou-4"
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional, String) The available zone that the CBS instance locates at.
* `project_id` - (Optional, String) ID of the project within the snapshot.
* `result_output_file` - (Optional, String) Used to save results.
* `snapshot_id` - (Optional, String) ID of the snapshot to be queried.
* `snapshot_name` - (Optional, String) Name of the snapshot to be queried.
* `storage_id` - (Optional, String) ID of the the CBS which this snapshot created from.
* `storage_usage` - (Optional, String) Types of CBS which this snapshot created from, and available values include `SYSTEM_DISK` and `DATA_DISK`.

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


