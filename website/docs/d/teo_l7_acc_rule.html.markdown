---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_l7_acc_rule"
sidebar_current: "docs-tencentcloud-datasource-teo_l7_acc_rule"
description: |-
  Use this data source to query detailed information of teo l7 acceleration rules
---

# tencentcloud_teo_l7_acc_rule

Use this data source to query detailed information of teo l7 acceleration rules

## Example Usage

```hcl
data "tencentcloud_teo_l7_acc_rule" "teo_l7_acc_rule" {
  zone_id = "zone-2qtuhspy7cr6"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Required, String) Zone ID.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `total_count` - (Computed, Int) Total count of L7 acceleration rules.
* `rules` - Rules content.

