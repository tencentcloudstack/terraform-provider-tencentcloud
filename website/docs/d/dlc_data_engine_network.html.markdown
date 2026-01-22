---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_data_engine_network"
sidebar_current: "docs-tencentcloud-datasource-dlc_data_engine_network"
description: |-
  Use this data source to query detailed information of DLC data engine network
---

# tencentcloud_dlc_data_engine_network

Use this data source to query detailed information of DLC data engine network

## Example Usage

```hcl
data "tencentcloud_dlc_data_engine_network" "example" {
  sort_by = "create-time"
  sorting = "desc"
  filters {
    name   = "engine-network-id"
    values = ["DataEngine_Network-g1sxyw8v"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - (Optional, List) Filter conditions are optional, engine-network-id--engine network ID, engine-network-state--engine network status.
* `result_output_file` - (Optional, String) Used to save results.
* `sort_by` - (Optional, String) Sort Field.
* `sorting` - (Optional, String) Sort order, asc or desc.

The `filters` object supports the following:

* `name` - (Required, String) Attribute name, if there are multiple filters, the relationship between filters is a logical OR relationship.
* `values` - (Required, List) Attribute value, if there are multiple values, the relationship between values is a logical OR relationship.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `engine_networks_infos` - Engine network information list.


