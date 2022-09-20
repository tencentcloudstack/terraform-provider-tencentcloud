---
subcategory: "TencentCloud EdgeOne(TEO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_teo_zone_available_plans"
sidebar_current: "docs-tencentcloud-datasource-teo_zone_available_plans"
description: |-
  Use this data source to query detailed information of teo zoneAvailablePlans
---

# tencentcloud_teo_zone_available_plans

Use this data source to query detailed information of teo zoneAvailablePlans

## Example Usage

```hcl
data "tencentcloud_teo_zone_available_plans" "zoneAvailablePlans" {
}
```

## Argument Reference

The following arguments are supported:

* `result_output_file` - (Optional, String) Used to save results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `plan_info_list` - Zone plans which current account can use.
  * `area` - Acceleration area of the plan. Valid value: `mainland`, `overseas`.
  * `currency` - Settlement Currency Type. Valid values: `CNY`, `USD`.
  * `flux` - The number of fluxes included in the zone plan. Unit: Byte.
  * `frequency` - Billing cycle. Valid values:- `y`: Billed by the year.- `m`: Billed by the month.- `h`: Billed by the hour.- `M`: Billed by the minute.- `s`: Billed by the second.
  * `plan_type` - Plan type.
  * `price` - Price of the plan. Unit: cent.
  * `request` - The number of requests included in the zone plan.
  * `site_number` - The number of zones this zone plan can bind.


