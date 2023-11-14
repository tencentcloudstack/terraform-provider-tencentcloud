package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
  month = &lt;nil&gt;
  end_month = &lt;nil&gt;
  member_uins = &lt;nil&gt;
  product_codes = &lt;nil&gt;
    }

`
