---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_standard_engine_resource_group_config_info"
sidebar_current: "docs-tencentcloud-resource-dlc_standard_engine_resource_group_config_info"
description: |-
  Provides a resource to create a DLC standard engine resource group config info
---

# tencentcloud_dlc_standard_engine_resource_group_config_info

Provides a resource to create a DLC standard engine resource group config info

~> **NOTE:** This resource must exclusive in one engine resource group, do not declare additional config resources of this conf context elsewhere.

~> **NOTE:** If you use the `tencentcloud_dlc_standard_engine_resource_group_config_info`. Please do not set `static_config_pairs` or `dynamic_config_pairs` in resource `tencentcloud_dlc_standard_engine_resource_group` simultaneously.

## Example Usage

```hcl
resource "tencentcloud_dlc_standard_engine_resource_group_config_info" "example" {
  engine_resource_group_name = "tf-example"
  static_conf_context {
    params {
      config_item  = "item1"
      config_value = "value1"
    }

    params {
      config_item  = "item2"
      config_value = "value2"
    }
  }

  dynamic_conf_context {
    params {
      config_item  = "item3"
      config_value = "value3"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `engine_resource_group_name` - (Required, String, ForceNew) Standard engine resource group name.
* `dynamic_conf_context` - (Optional, List) Dynamic config context.
* `static_conf_context` - (Optional, List) Static config context.

The `dynamic_conf_context` object supports the following:

* `params` - (Optional, Set) Collection of bound working groups.

The `params` object of `dynamic_conf_context` supports the following:

* `config_item` - (Required, String) Configuration item.
* `config_value` - (Required, String) Configuration value.

The `params` object of `static_conf_context` supports the following:

* `config_item` - (Required, String) Configuration item.
* `config_value` - (Required, String) Configuration value.

The `static_conf_context` object supports the following:

* `params` - (Optional, Set) Collection of bound working groups.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

DLC standard engine resource group config info can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_standard_engine_resource_group_config_info.example tf-example
```

