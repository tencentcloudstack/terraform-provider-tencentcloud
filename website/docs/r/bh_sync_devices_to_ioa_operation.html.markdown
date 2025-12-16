---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_sync_devices_to_ioa_operation"
sidebar_current: "docs-tencentcloud-resource-bh_sync_devices_to_ioa_operation"
description: |-
  Provides a resource to create a BH sync devices to ioa operation
---

# tencentcloud_bh_sync_devices_to_ioa_operation

Provides a resource to create a BH sync devices to ioa operation

## Example Usage

```hcl
resource "tencentcloud_bh_sync_devices_to_ioa_operation" "example" {
  device_id_set = [
    1934,
    1964,
    1895,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `device_id_set` - (Required, Set: [`Int`], ForceNew) Asset ID collection. Assets must be bound to bastion host instances that support IOA functionality. Maximum 200 assets can be synchronized at a time.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



