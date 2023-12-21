package eb_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudEbPlatformProductsDataSource_basic -v
func TestAccTencentCloudEbPlatformProductsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccStepSetRegion(t, "ap-chongqing")
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccEbPlatformProductsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_eb_platform_products.platform_products"),
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
