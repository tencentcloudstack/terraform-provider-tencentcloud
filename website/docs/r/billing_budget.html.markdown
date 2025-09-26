---
subcategory: "Billing"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_billing_budget"
sidebar_current: "docs-tencentcloud-resource-billing_budget"
description: |-
  Provides a resource to create a billing billing_budget
---

# tencentcloud_billing_budget

Provides a resource to create a billing billing_budget

## Example Usage

```hcl
resource "tencentcloud_billing_budget" "billing_budget" {
  budget_name  = "tf-test"
  cycle_type   = "MONTH"
  period_begin = "2025-09"
  period_end   = "2026-09"
  plan_type    = "FIX"
  budget_quota = "10000.00"
  bill_type    = "BILL"
  fee_type     = "REAL_COST"

  budget_note = "budget_note"

  warn_json {
    warn_type       = "ACTUAL"
    cal_type        = "PERCENTAGE"
    threshold_value = "60"
  }
  dimensions_range {
    business      = ["p_cvm"]
    pay_mode      = ["prePay"]
    product_codes = ["sp_cvm_s6"]
    zone_ids      = ["100006"]
    region_ids    = ["1"]
    project_ids   = ["0"]
    action_types  = ["prepay_purchase"]
  }
  wave_threshold_json {
    warn_type   = "ACTUAL"
    threshold   = "20"
    meta_type   = "chain"
    period_type = "day"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bill_type` - (Required, String) BILL: system bill, CONSUMPTION: consumption bill.
* `budget_name` - (Required, String) Budget name.
* `budget_quota` - (Required, String) Budget value limit. Transfer fixed value when the budget plan type is FIX(Fixed Budget); Passed when the budget plan type is CYCLE(Planned Budget)[{"dateDesc":"2025-07","quota":"1000"},{"dateDesc":"2025-08","quota":"2000"}].
* `cycle_type` - (Required, String) Cycle type, valid values: DAY, MONTH, QUARTER, YEAR.
* `fee_type` - (Required, String) COST original price, REAL_COST actual cost, CASH cash, INCENTIVE gift, VOUCHER voucher, TRANSFER share, TAX tax, AMOUNT_BEFORE_TAX cash payment (before tax).
* `period_begin` - (Required, String) Valid period starting time 2025-01-01(cycle: days) / 2025-01 (cycle: months).
* `period_end` - (Required, String) Expiration period end time 2025-12-01(cycle: days) / 2025-12 (cycle: months).
* `plan_type` - (Required, String) FIX: fixed budget, CYCLE: planned budget.
* `warn_json` - (Required, List) Threshold reminder.
* `budget_note` - (Optional, String) Budget remarks.
* `dimensions_range` - (Optional, List) Budget dimension range conditions.
* `wave_threshold_json` - (Optional, List) Volatility reminder.

The `dimensions_range` object supports the following:

* `action_types` - (Optional, Set) Action types.
* `business` - (Optional, Set) Products.
* `component_codes` - (Optional, Set) Component codes.
* `consumption_types` - (Optional, Set) Consumption types.
* `owner_uins` - (Optional, Set) Owner uins.
* `pay_mode` - (Optional, Set) Pay mode.
* `payer_uins` - (Optional, Set) Payer uins.
* `product_codes` - (Optional, Set) Sub-product.
* `project_ids` - (Optional, Set) Project ids.
* `region_ids` - (Optional, Set) Region ids.
* `tags` - (Optional, List) Tags.
* `tree_node_uniq_keys` - (Optional, Set) Unique key for end-level ledger unit.
* `zone_ids` - (Optional, Set) Zone ids.

The `tags` object of `dimensions_range` supports the following:

* `tag_key` - (Optional, String) Tag key.
* `tag_value` - (Optional, Set) Tag value.

The `warn_json` object supports the following:

* `cal_type` - (Required, String) PERCENTAGE: Percentage of budget amount, ABS: fixed value.
* `threshold_value` - (Required, String) Threshold (greater than or equal to 0).
* `warn_type` - (Required, String) ACTUAL: actual amount, FORECAST: forecast amount.

The `wave_threshold_json` object supports the following:

* `meta_type` - (Optional, String) Alarm type: chain month-on-month, yoy year-on-year, fix fixed value
 (Supported types: daily month-on-month chain day, daily month-on-year chain weekday, daily month-on-year monthly month-on-year fixed value fix day, month-on-month chain month, monthly fixed value fix month).
* `period_type` - (Optional, String) Alarm dimension: day day, month month, weekday week
 (Support types: day-to-day chain day, day-to-year weekly dimension chain weekday, day-to-year monthly dimension yoy day, daily fixed value fix day, month-to-month chain month, monthly fixed value fix month).
* `threshold` - (Optional, String) Volatility threshold (greater than or equal to 0).
* `warn_type` - (Optional, String) ACTUAL: actual amount, FORECAST: forecast amount.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

billing billing_budget can be imported using the id, e.g.

```
terraform import tencentcloud_billing_budget.billing_budget billing_budget_id
```

