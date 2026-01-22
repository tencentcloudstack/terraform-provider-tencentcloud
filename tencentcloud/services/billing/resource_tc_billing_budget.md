Provides a resource to create a billing billing_budget

Example Usage

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

Import

billing billing_budget can be imported using the id, e.g.

```
terraform import tencentcloud_billing_budget.billing_budget billing_budget_id
```
