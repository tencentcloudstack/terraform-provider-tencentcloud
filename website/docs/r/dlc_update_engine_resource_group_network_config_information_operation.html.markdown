---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_update_engine_resource_group_network_config_information_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_update_engine_resource_group_network_config_information_operation"
description: |-
  Provides a resource to create a DLC update engine resource group network config information operation
---

# tencentcloud_dlc_update_engine_resource_group_network_config_information_operation

Provides a resource to create a DLC update engine resource group network config information operation

## Example Usage

```hcl
resource "tencentcloud_dlc_update_engine_resource_group_network_config_information_operation" "example" {
  engine_resource_group_id = "rg-b6fxxxxxx2a0"
}
```

## Argument Reference

The following arguments are supported:

* `engine_resource_group_id` - (Required, String, ForceNew) Engine resource group ID.
* `network_config_names` - (Optional, Set: [`String`], ForceNew) A collection of network configuration names bound to the resource group.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



