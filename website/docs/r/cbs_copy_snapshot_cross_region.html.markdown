---
subcategory: "Cloud Block Storage(CBS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cbs_copy_snapshot_cross_region"
sidebar_current: "docs-tencentcloud-resource-cbs_copy_snapshot_cross_region"
description: |-
  Provides a resource to create a CBS copy snapshot cross region resource.
---

# tencentcloud_cbs_copy_snapshot_cross_region

Provides a resource to create a CBS copy snapshot cross region resource.

## Example Usage

```hcl
resource "tencentcloud_cbs_copy_snapshot_cross_region" "example" {
  snapshot_id         = "snap-xxxxxxxx"
  destination_regions = ["ap-shanghai", "ap-chengdu"]

  snapshot_name = "my-copied-snapshot"
}
```

## Argument Reference

The following arguments are supported:

* `destination_regions` - (Required, List: [`String`], ForceNew) Target region names for cross-region copy.
* `snapshot_id` - (Required, String, ForceNew) Source snapshot ID for cross-region copy.
* `delete_bind_images` - (Optional, Bool, ForceNew) Whether to force-delete images associated with the snapshots when deleting.
* `snapshot_name` - (Optional, String, ForceNew) Name of the copied snapshot.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `snapshot_copy_result_set` - Cross-region copy results.
  * `code` - Error code, Success on success.
  * `destination_region` - Destination region for cross-region copy.
  * `message` - Error message, empty string on success.
  * `snapshot_id` - New snapshot ID in the destination region.


## Import

CBS copy snapshot cross region can be imported using the composite ID, e.g.

```
$ terraform import tencentcloud_cbs_copy_snapshot_cross_region.example snap-xxxxxxxx#snap-yyyyyyyy
```

