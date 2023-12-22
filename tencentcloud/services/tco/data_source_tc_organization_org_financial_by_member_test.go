package tco_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudOrganizationOrgFinancialByMemberDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgFinancialByMemberDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_organization_org_financial_by_member.org_financial_by_member")),
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
