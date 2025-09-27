---
subcategory: "Data Lake Compute(DLC)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dlc_data_mask_strategy"
sidebar_current: "docs-tencentcloud-resource-dlc_data_mask_strategy"
description: |-
  Provides a resource to create a DLC data mask strategy
---

# tencentcloud_dlc_data_mask_strategy

Provides a resource to create a DLC data mask strategy

## Example Usage

```hcl
resource "tencentcloud_dlc_data_mask_strategy" "example" {
  strategy {
    strategy_name = "tf-example"
    strategy_desc = "description."
    groups {
      work_group_id = 70220
      strategy_type = "MASK"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `strategy` - (Optional, List) The data masking strategy.

The `groups` object of `strategy` supports the following:

* `strategy_type` - (Optional, String) The type of the data masking strategy. Supported value: MASK/MASK_NONE/MASK_NULL/MASK_HASH/MASK_SHOW_LAST_4/MASK_SHOW_FIRST_4/MASK_DATE_SHOW_YEAR.
* `work_group_id` - (Optional, Int) The unique ID of the work group.

The `strategy` object supports the following:

* `groups` - (Optional, List) Collection of bound working groups.
* `strategy_desc` - (Optional, String) The description of the data masking strategy.
* `strategy_name` - (Optional, String) The name of the data masking strategy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

DLC data mask strategy can be imported using the id, e.g.

```
terraform import tencentcloud_dlc_data_mask_strategy.example 2fcab650-11a8-44ef-bf58-19c22af601b6
```

