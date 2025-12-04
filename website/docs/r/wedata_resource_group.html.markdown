---
subcategory: "Wedata"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_wedata_resource_group"
sidebar_current: "docs-tencentcloud-resource-wedata_resource_group"
description: |-
  Provides a resource to create a WeData resource group
---

# tencentcloud_wedata_resource_group

Provides a resource to create a WeData resource group

## Example Usage

```hcl
resource "tencentcloud_wedata_resource_group" "example" {
  name = "tf_example"
  type {
    resource_group_type = "Integration"
    integration {
      real_time_data_sync {
        specification = "i32c"
        number        = 1
      }

      offline_data_sync {
        specification = "integrated"
        number        = 2
      }
    }
  }

  auto_renew_enabled = false
  purchase_period    = 1
  vpc_id             = "vpc-ds5rpnxh"
  subnet             = "subnet-fz7rw5zq"
  resource_region    = "ap-beijing-fsi"
}
```

## Argument Reference

The following arguments are supported:

* `auto_renew_enabled` - (Required, Bool) Whether auto-renewal is enabled.
* `name` - (Required, String, ForceNew) Resource group name. The name for creating a general resource group must start with a letter, can contain letters, numbers, underscores (_), and up to 64 characters.
* `purchase_period` - (Required, Int) Purchase duration, in months.
* `resource_region` - (Required, String, ForceNew) Resource purchase region.
* `subnet` - (Required, String, ForceNew) Subnet.
* `type` - (Required, List, ForceNew) Information about the activated resource group.
* `vpc_id` - (Required, String, ForceNew) VPC ID.
* `associated_project_id` - (Optional, String, ForceNew) Associated project space project ID.

The `data_service` object of `type` supports the following:

* `number` - (Required, Int) Quantity.
* `specification` - (Required, String) Resource group specification.

The `integration` object of `type` supports the following:

* `offline_data_sync` - (Optional, List) Offline integration resource group.

- integrated (Offline data synchronization - 8C16G)
- i16 (Offline data synchronization - 8C32G).
* `real_time_data_sync` - (Optional, List) Real-time integration resource group.

- i32c (Real-time data synchronization - 16C64G).

The `offline_data_sync` object of `integration` supports the following:

* `number` - (Required, Int) Quantity.
* `specification` - (Required, String) Resource group specification.

The `real_time_data_sync` object of `integration` supports the following:

* `number` - (Required, Int) Quantity.
* `specification` - (Required, String) Resource group specification.

The `schedule` object of `type` supports the following:

* `number` - (Required, Int) Quantity.
* `specification` - (Required, String) Resource group specification.

The `type` object supports the following:

* `resource_group_type` - (Required, String) Resource group type.

- Schedule --- Scheduling resource group
- Integration --- Integration resource group  
- DataService -- Data service resource group.
* `data_service` - (Optional, List) Data service resource group (Integration, scheduling, and data service resource groups cannot be purchased simultaneously).

- ds_t (Test specification)
- ds_s (Basic specification)
- ds_m (Popular specification)
- ds_l (Professional specification).
* `integration` - (Optional, List) Integration resource group, subdivided into real-time resource group and offline resource group (Integration, scheduling, and data service resource groups cannot be purchased simultaneously).
* `schedule` - (Optional, List) Scheduling resource group (Integration, scheduling, and data service resource groups cannot be purchased simultaneously).

- s_test (Test specification)
- s_small (Basic specification)
- s_medium (Popular specification)
- s_large (Professional specification).

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `resource_group_id` - Resource group ID.


