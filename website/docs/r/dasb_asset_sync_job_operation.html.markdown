---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_asset_sync_job_operation"
sidebar_current: "docs-tencentcloud-resource-dasb_asset_sync_job_operation"
description: |-
  Provides a resource to create a dasb asset sync job
---

# tencentcloud_dasb_asset_sync_job_operation

Provides a resource to create a dasb asset sync job

## Example Usage

```hcl
resource "tencentcloud_dasb_asset_sync_job_operation" "example" {
  category = 1
}
```

## Argument Reference

The following arguments are supported:

* `category` - (Required, Int, ForceNew) Synchronize asset categories, 1- Host assets, 2- Database assets.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateAssetSyncJob](https://cloud.tencent.com/document/api/1025/86957)
* Create: [DescribeAssetSyncStatus](https://cloud.tencent.com/document/api/1025/86953)

