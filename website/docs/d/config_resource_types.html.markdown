---
subcategory: "Config"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_config_resource_types"
sidebar_current: "docs-tencentcloud-datasource-config_resource_types"
description: |-
  Use this data source to query the list of resource types supported by Tencent Cloud Config.
---

# tencentcloud_config_resource_types

Use this data source to query the list of resource types supported by Tencent Cloud Config.

## Example Usage

```hcl
data "tencentcloud_config_resource_types" "example" {}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `resource_type_list` - Supported resource type list.
  * `product_name` - Product name.
  * `product` - Product code (e.g. CAM).
  * `resource_type_name` - Resource type name.
  * `resource_type` - Resource type identifier (e.g. QCS::CAM::Group).


