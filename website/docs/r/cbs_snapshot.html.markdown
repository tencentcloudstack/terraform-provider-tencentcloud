---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_snapshot"
sidebar_current: "docs-tencentcloud-resource-cbs_snapshot"
description: |-
  Provides a resource to create a CBS snapshot.
---

# tencentcloud_cbs_snapshot

Provides a resource to create a CBS snapshot.

## Example Usage

```hcl
resource "tencentcloud_cbs_snapshot" "example" {
  storage_id    = "disk-1i9gxxi8"
  snapshot_name = "tf-example"
  disk_usage    = "DATA_DISK"
  tags = {
    createBy = "Terraform"
  }
}
```

## Argument Reference

The following arguments are supported:

* `snapshot_name` - (Required, String) Name of the snapshot.
* `storage_id` - (Required, String, ForceNew) ID of the the CBS which this snapshot created from.
* `disk_usage` - (Optional, String, ForceNew) The type of cloud disk associated with the snapshot: SYSTEM_DISK: system disk; DATA_DISK: data disk. If not filled in, the snapshot type will be consistent with the cloud disk type. This parameter is used in some scenarios where users need to create a data disk snapshot from the system disk for shared use.
* `tags` - (Optional, Map) The available tags within this CBS Snapshot.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of snapshot.
* `disk_type` - (**Deprecated**) It has been deprecated from version 1.82.14. Please use `disk_usage` instead. Types of CBS which this snapshot created from.
* `percent` - Snapshot creation progress percentage. If the snapshot has created successfully, the constant value is 100.
* `snapshot_status` - Status of the snapshot.
* `storage_size` - Volume of storage which this snapshot created from.


## Import

CBS snapshot can be imported using the id, e.g.

```
$ terraform import tencentcloud_cbs_snapshot.example snap-3sa3f39b
```

