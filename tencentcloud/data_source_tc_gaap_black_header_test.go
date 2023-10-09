package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapBlackHeaderDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapBlackHeaderDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_black_header.black_header"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_black_header.black_header", "black_headers.#"),
				),
			},
		},
	})
}

const testAccGaapBlackHeaderDataSource = `

data "tencentcloud_gaap_black_header" "black_header" {
}
`
