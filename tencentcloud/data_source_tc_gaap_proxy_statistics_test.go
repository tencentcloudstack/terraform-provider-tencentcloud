package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapProxyStatisticsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapProxyStatisticsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_proxy_statistics.proxy_statistics"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_proxy_statistics.proxy_statistics", "statistics_data.#"),
				),
			},
		},
	})
}

const testAccGaapProxyStatisticsDataSource = `
data "tencentcloud_gaap_proxy_statistics" "proxy_statistics" {
	proxy_id = "link-8lpyo88p"
	start_time = "2023-10-09 00:00:00"
	end_time = "2023-10-09 23:59:59"
	metric_names = ["InBandwidth", "OutBandwidth", "InFlow", "OutFlow", "InPackets", "OutPackets", "Concurrent", "HttpQPS", "HttpsQPS", "Latency", "PacketLoss"]
	granularity = 300
}
`
