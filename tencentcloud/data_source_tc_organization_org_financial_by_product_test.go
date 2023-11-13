package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudOrganizationOrgFinancialByProductDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOrganizationOrgFinancialByProductDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_organization_org_financial_by_product.org_financial_by_product")),
			},
		},
	})
}

const testAccOrganizationOrgFinancialByProductDataSource = `

data "tencentcloud_organization_org_financial_by_product" "org_financial_by_product" {
  month = &lt;nil&gt;
  end_month = &lt;nil&gt;
  member_uins = &lt;nil&gt;
  product_codes = &lt;nil&gt;
    }

`
