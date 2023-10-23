package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapProxyAndStatisticsListenersDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyAndStatisticsListenersDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxy_and_statistics_listeners.proxy_and_statistics_listeners"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxy_and_statistics_listeners.proxy_and_statistics_listeners", "proxy_set.#"),
				),
			},
		},
	})
}

const testAccGaapProxyAndStatisticsListenersDataSource = `
data "tencentcloud_gaap_proxy_and_statistics_listeners" "proxy_and_statistics_listeners" {
	project_id = 0
}
`
