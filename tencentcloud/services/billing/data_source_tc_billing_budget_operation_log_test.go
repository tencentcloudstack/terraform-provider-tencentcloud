package billing_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudBillingBudgetOperationLogDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccBillingBudgetOperationLogDataSource,
			Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_billing_budget_operation_log.billing_budget_operation_log")),
		}},
	})
}

const testAccBillingBudgetOperationLogDataSource = `
data "tencentcloud_billing_budget_operation_log" "billing_budget_operation_log" {
    budget_id = "1971489821259956225"
}
`
