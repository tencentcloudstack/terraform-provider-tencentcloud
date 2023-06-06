---
subcategory: "Auto Scaling(AS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_as_advices"
sidebar_current: "docs-tencentcloud-datasource-as_advices"
description: |-
  Use this data source to query detailed information of as advices
---

# tencentcloud_as_advices

Use this data source to query detailed information of as advices

## Example Usage

```hcl
data "tencentcloud_as_advices" "advices" {
  auto_scaling_group_ids = ["asc-lo0b94oy"]
}
```

## Argument Reference

The following arguments are supported:

* `auto_scaling_group_ids` - (Required, Set: [`String`]) List of scaling groups to be queried. Upper limit: 100.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `auto_scaling_advice_set` - A collection of suggestions for scaling group configurations.
  * `advices` - A collection of suggestions for scaling group configurations.
    * `detail` - Problem Details.
    * `problem` - Problem Description.
    * `solution` - Recommended resolutions.
  * `auto_scaling_group_id` - Auto scaling group ID.
  * `level` - Scaling group warning level. Valid values: NORMAL, WARNING, CRITICAL.


