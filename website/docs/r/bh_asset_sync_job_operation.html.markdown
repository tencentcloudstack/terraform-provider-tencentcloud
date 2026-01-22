---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_asset_sync_job_operation"
sidebar_current: "docs-tencentcloud-resource-bh_asset_sync_job_operation"
description: |-
  Provides a resource to create a BH asset sync job operation
---

# tencentcloud_bh_asset_sync_job_operation

Provides a resource to create a BH asset sync job operation

## Example Usage

```hcl
resource "tencentcloud_bh_asset_sync_job_operation" "example" {
  category = 1
}
```

## Argument Reference

The following arguments are supported:

* `category` - (Required, Int, ForceNew) Asset synchronization category. 1 - host assets, 2 - database assets, 3 - Container assets.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



