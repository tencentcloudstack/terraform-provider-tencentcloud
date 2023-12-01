package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapProxyDetailDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyDetailDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxy_detail.proxy_detail"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxy_detail.proxy_detail", "proxy_detail.#"),
				),
			},
		},
	})
}

const testAccGaapProxyDetailDataSource = `
data "tencentcloud_gaap_proxy_detail" "proxy_detail" {
	proxy_id = "link-8lpyo88p"
}
`
