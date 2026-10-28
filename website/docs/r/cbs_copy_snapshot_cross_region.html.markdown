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
  snapshot_id        = "snap-07ttd84z"
  destination_region = "ap-beijing"
  snapshot_name      = "tf-example-copy-snapshot"
}
```

## Argument Reference

The following arguments are supported:

* `destination_region` - (Required, String, ForceNew) Target region name for cross-region copy.
* `snapshot_id` - (Required, String, ForceNew) Source snapshot ID for cross-region copy.
* `delete_bind_images` - (Optional, Bool) Whether to force-delete images associated with the snapshots when deleting.
* `snapshot_name` - (Optional, String, ForceNew) Name of the copied snapshot.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to `15m`) Used when creating the resource.
* `delete` - (Defaults to `15m`) Used when deleting the resource.

