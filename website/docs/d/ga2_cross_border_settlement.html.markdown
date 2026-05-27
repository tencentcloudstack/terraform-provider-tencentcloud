---
subcategory: "Global Accelerator 2.0(GA2)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ga2_cross_border_settlement"
sidebar_current: "docs-tencentcloud-datasource-ga2_cross_border_settlement"
description: |-
  Use this data source to query GA2 (Global Accelerator 2.0) cross-border settlement traffic usage information.
---

# tencentcloud_ga2_cross_border_settlement

Use this data source to query GA2 (Global Accelerator 2.0) cross-border settlement traffic usage information.

## Example Usage

### Query cross-border settlement traffic

```hcl
data "tencentcloud_ga2_cross_border_settlement" "example" {
  global_accelerator_id = "ga2-xxxxxxxx"
  accelerate_region     = "ap-guangzhou"
  endpoint_group_region = "ap-singapore"
  settlement_month      = 202501
}
```

## Argument Reference

The following arguments are supported:

* `accelerate_region` - (Required, String) Acceleration region.
* `endpoint_group_region` - (Required, String) Endpoint group region.
* `global_accelerator_id` - (Required, String) Global accelerator instance ID.
* `settlement_month` - (Required, Int) Billing year-month time.
* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `traffic` - Traffic usage in GB with 6 decimal places precision.


