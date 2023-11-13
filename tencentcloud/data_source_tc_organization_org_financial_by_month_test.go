package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudOrganizationOrgFinancialByMonthDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgFinancialByMonthDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_organization_org_financial_by_month.org_financial_by_month")),
			},
		},
	})
}

const testAccOrganizationOrgFinancialByMonthDataSource = `

data "tencentcloud_organization_org_financial_by_month" "org_financial_by_month" {
  end_month = "2021-01"
  member_uins = &lt;nil&gt;
  product_codes = &lt;nil&gt;
  }

`
