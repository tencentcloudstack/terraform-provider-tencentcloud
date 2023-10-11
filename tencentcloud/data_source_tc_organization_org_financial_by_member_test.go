package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOrganizationOrgFinancialByMemberDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgFinancialByMemberDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_organization_org_financial_by_member.org_financial_by_member")),
			},
		},
	})
}

const testAccOrganizationOrgFinancialByMemberDataSource = `

data "tencentcloud_organization_org_financial_by_member" "org_financial_by_member" {
  month = "2023-05"
  end_month = "2023-10"
  member_uins = [100015591986,100029796005]
    }

`
