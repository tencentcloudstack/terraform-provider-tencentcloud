package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudEbPlatformProductsDataSource_basic -v
func TestAccTencentCloudEbPlatformProductsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-chongqing")
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbPlatformProductsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_eb_platform_products.platform_products"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_eb_platform_products.platform_products", "platform_products.#"),
					resource.TestCheckResourceAttr("data.tencentcloud_eb_platform_products.platform_products", "platform_products.0.product_type", "eb_platform_test"),
				),
			},
		},
	})
}

const testAccEbPlatformProductsDataSource = `

data "tencentcloud_eb_platform_products" "platform_products" {
}

`
