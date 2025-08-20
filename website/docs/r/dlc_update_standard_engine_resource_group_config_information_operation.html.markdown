---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_update_standard_engine_resource_group_config_information_operation"
sidebar_current: "docs-tencentcloud-resource-dlc_update_standard_engine_resource_group_config_information_operation"
description: |-
  Provides a resource to create a DLC update standard engine resource group config information operation
---

# tencentcloud_dlc_update_standard_engine_resource_group_config_information_operation

Provides a resource to create a DLC update standard engine resource group config information operation

## Example Usage

```hcl
resource "tencentcloud_dlc_update_standard_engine_resource_group_config_information_operation" "example" {
  engine_resource_group_name = "tf-example"
  update_conf_context {
    config_type = "StaticConfigType"
    params {
      config_item  = "spark.sql.shuffle.partitions"
      config_value = "300"
      operate      = "ADD"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `engine_resource_group_name` - (Required, String, ForceNew) Engine resource group name.
* `update_conf_context` - (Required, List, ForceNew) Configuration that needs to be updated.

The `params` object of `update_conf_context` supports the following:

* `config_item` - (Optional, String, ForceNew) Parameter key.
* `config_value` - (Optional, String, ForceNew) Parameter value.
* `operate` - (Optional, String, ForceNew) Send operations, support: ADD, DELETE, MODIFY.

The `update_conf_context` object supports the following:

* `config_type` - (Required, String, ForceNew) Parameter type, optional: StaticConfigType, DynamicConfigType.
* `params` - (Required, List) Configuration array of parameters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



