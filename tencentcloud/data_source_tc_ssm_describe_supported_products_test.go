package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSsmDescribeSupportedProductsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSsmDescribeSupportedProductsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ssm_describe_supported_products.describe_supported_products")),
			},
		},
	})
}

const testAccSsmDescribeSupportedProductsDataSource = `

data "tencentcloud_ssm_describe_supported_products" "describe_supported_products" {
  }

`
