package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Check: resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_organization_org_financial_by_product.org_financial_by_product"),
					resource.TestCheckResourceAttr("data.tencentcloud_organization_org_financial_by_product.org_financial_by_product", "month", "2023-05"),
					resource.TestCheckResourceAttr("data.tencentcloud_organization_org_financial_by_product.org_financial_by_product", "end_month", "2023-09"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_organization_org_financial_by_product.org_financial_by_product", "product_codes.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_organization_org_financial_by_product.org_financial_by_product", "product_codes.0", "p_eip"),
				),
			},
		},
	})
}

const testAccOrganizationOrgFinancialByProductDataSource = `

data "tencentcloud_organization_org_financial_by_product" "org_financial_by_product" {
  month = "2023-05"
  end_month = "2023-09"
  product_codes = ["p_eip"]
}

`
