package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapProxiesStatusDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxiesStatusDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxies_status.proxies_status"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxies_status.proxies_status", "instance_status_set.#"),
				),
			},
		},
	})
}

const testAccGaapProxiesStatusDataSource = `
data "tencentcloud_gaap_proxies_status" "proxies_status" {
  proxy_ids = ["link-8lpyo88p"]
}
`
