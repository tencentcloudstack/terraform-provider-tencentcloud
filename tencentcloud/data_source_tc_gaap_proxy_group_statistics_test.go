package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapProxyGroupStatisticsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyGroupStatisticsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxy_group_statistics.proxy_group_statistics"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxy_group_statistics.proxy_group_statistics", "statistics_data.#"),
				),
			},
		},
	})
}

const testAccGaapProxyGroupStatisticsDataSource = `
data "tencentcloud_gaap_proxy_group_statistics" "proxy_group_statistics" {
	group_id = "link-8lpyo88p"
	start_time = "2023-10-09 00:00:00"
	end_time = "2023-10-09 23:59:59"
	metric_names = ["InBandwidth", "OutBandwidth", "InFlow", "OutFlow"]
	granularity = 300
}
`
