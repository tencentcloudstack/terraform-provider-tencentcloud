---
layout: "tencentcloud"
page_title: "tencentcloud: tencentcloud_snapshot"
sidebar_current: "docs-tencentcloud-resource-snapshot"
description: |-
  Provides a snapshot resource.
---

# tencentcloud_snapshot

Provides a snapshot resource.

## Example Usage

```hcl
resource "tencentcloud_snapshot" "my-snapshot" {
  storage_id	= "disk-4vmyor8k"
  snapshot_name = "my-snapshot"
}
```

## Argument Reference

The following arguments are supported:

* `storage_id` - (Required) Source Storage to create this snapshot.
* `snapshot_name` - (Optional) The name of the snapshot. This snapshot_name can have a string of 1 to 64 characters. It is supported to modify `snapshot_name` after the snapshot is created.


## Attributes Reference

The following attributes are exported:

* `id` - The snapshot ID, something looks like `snapshot-xxxxxx`.
* `storage_id` - The storage ID which this snapshot created from.
* `storage_size` - The size of assiciated storage. You can create new larger or same size storage use this snapshot.
* `snapshot_name` - Name of snapshot
* `percent` - The creation progress of this snapshot.
* `disk_type` - The disk type of this snapshot, `root` or `data`.
* `snapshot_status` - The status of this snapshot. "creating" means the snapshot is creating; "normal" means the snapshot is ready to use.
