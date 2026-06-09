---
subcategory: "Cloud Virtual Machine(CVM)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_cvm_resource_pool_pack"
sidebar_current: "docs-tencentcloud-resource-cvm_resource_pool_pack"
description: |-
  Provides a resource to create CVM resource pool packs for advance compute resource pool management.
---

# tencentcloud_cvm_resource_pool_pack

Provides a resource to create CVM resource pool packs for advance compute resource pool management.

Resource pool packs allow users to purchase and manage compute resource pools in advance. This resource supports standard Terraform operations: create, read, and delete. Update operations are not available from CVM API, all field changes will trigger a destroy and recreate.

## Example Usage

### Basic Usage

```hcl
resource "tencentcloud_cvm_resource_pool_pack" "pool_pack" {
  zone                              = "ap-guangzhou-7"
  instance_type                     = "SA9.96XLARGE1152"
  period                            = 12
  resource_pool_pack_type           = "EXCLUSIVE"
  auto_placement                    = true
  dedicated_resource_pool_pack_name = "my-resource-pool"
  renew_flag                        = "NOTIFY_AND_MANUAL_RENEW"
}
```

### With Auto Renewal

```hcl
resource "tencentcloud_cvm_resource_pool_pack" "pool_pack_auto_renew" {
  zone                              = "ap-guangzhou-7"
  instance_type                     = "SA9.96XLARGE1152"
  period                            = 6
  resource_pool_pack_type           = "EXCLUSIVE"
  auto_placement                    = true
  dedicated_resource_pool_pack_name = "auto-renew-pool"
  renew_flag                        = "NOTIFY_AND_AUTO_RENEW"
}
```

## Argument Reference

The following arguments are supported:

* `instance_type` - (Required, String, ForceNew) Instance type for the resource pool pack. Only half-machine/full-machine specifications are supported. Format: SA9.96XLARGE1152 (SA9 half-machine).
* `period` - (Required, Int, ForceNew) The period of the resource pool pack in months. Range: 1-60.
* `zone` - (Required, String, ForceNew) The availability zone where the resource pool pack is located. Format: ap-guangzhou-6.
* `auto_placement` - (Optional, Bool, ForceNew) Auto placement switch. Default: true. When enabled, the system will search for suitable pools in pools with this capability enabled when creating instances without specifying a resource pool.
* `dedicated_resource_pool_pack_name` - (Optional, String, ForceNew) The name of the resource pool pack. Length: 1-60 characters, supports Chinese, English, numbers, hyphens '-', and underscores '_'.
* `renew_flag` - (Optional, String, ForceNew) Auto renewal flag. Options: NOTIFY_AND_AUTO_RENEW (notify and auto renew), NOTIFY_AND_MANUAL_RENEW (notify and manual renew, default), DISABLE_NOTIFY_AND_MANUAL_RENEW (do not notify and manual renew).
* `resource_pool_pack_type` - (Optional, String, ForceNew) Resource pool pack type. Options: EXCLUSIVE (exclusive, default), SHARED (shared). Note: Only EXCLUSIVE is supported in the first phase.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `dedicated_resource_pack_id` - The ID of the created resource pool pack.
* `end_time` - Resource pool pack expiration time. Format: YYYY-MM-DDThh:mm:ssZ.
* `instance_family` - Instance family. Format: SA9.
* `start_time` - Resource pool pack creation time. Format: YYYY-MM-DDThh:mm:ssZ.
* `status` - Resource pool pack status. Values: CREATING (creating), ACTIVE (running), FAILED (creation failed), RETIRED (expired).


## Import

CVM resource pool packs can be imported using the resource pool pack ID, e.g.

```
terraform import tencentcloud_cvm_resource_pool_pack.pool_pack rpp-xxxxx
```

