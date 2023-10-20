package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapListenerStatisticsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapListenerStatisticsDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_listener_statistics.listener_statistics"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_gaap_listener_statistics.listener_statistics", "statistics_data.#"),
				),
			},
		},
	})
}

const testAccGaapListenerStatisticsDataSource = `
data "tencentcloud_gaap_listener_statistics" "listener_statistics" {
	listener_id = "listener-2s3lghkv"
	start_time = "2023-10-19 00:00:00"
	end_time = "2023-10-19 23:59:59"
	metric_names = ["InBandwidth", "OutBandwidth", "InPackets", "OutPackets", "Concurrent"]
	granularity = 300
}
`
