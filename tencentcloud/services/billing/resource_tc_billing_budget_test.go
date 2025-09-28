package billing_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBillingBudgetResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBillingBudget,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_billing_budget.billing_budget", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_billing_budget.billing_budget", "budget_name"),
					resource.TestCheckResourceAttrSet("tencentcloud_billing_budget.billing_budget", "cycle_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_billing_budget.billing_budget", "period_begin"),
					resource.TestCheckResourceAttrSet("tencentcloud_billing_budget.billing_budget", "period_end"),
					resource.TestCheckResourceAttrSet("tencentcloud_billing_budget.billing_budget", "plan_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_billing_budget.billing_budget", "budget_quota"),
					resource.TestCheckResourceAttrSet("tencentcloud_billing_budget.billing_budget", "bill_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_billing_budget.billing_budget", "fee_type"),
					resource.TestCheckResourceAttr("tencentcloud_billing_budget.billing_budget", "warn_json.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_billing_budget.billing_budget", "dimensions_range.#", "1"),
					resource.TestCheckResourceAttr("tencentcloud_billing_budget.billing_budget", "wave_threshold_json.#", "1"),
				),
			},
			{
				ResourceName:      "tencentcloud_billing_budget.billing_budget",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccBillingBudget = `
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
`
