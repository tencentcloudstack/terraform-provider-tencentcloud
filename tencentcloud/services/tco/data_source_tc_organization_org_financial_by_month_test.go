package tco_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOrganizationOrgFinancialByMonthDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgFinancialByMonthDataSource,
				Check: resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_organization_org_financial_by_month.org_financial_by_month"),
					resource.TestCheckResourceAttr("data.tencentcloud_organization_org_financial_by_month.org_financial_by_month", "end_month", "2023-05"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_organization_org_financial_by_month.org_financial_by_month", "member_uins.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_organization_org_financial_by_month.org_financial_by_month", "member_uins.0", "100026517717"),
				),
			},
		},
	})
}

const testAccOrganizationOrgFinancialByMonthDataSource = `

data "tencentcloud_organization_org_financial_by_month" "org_financial_by_month" {
  end_month = "2023-05"
  member_uins = [100026517717]
  }

`
