package gaap_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudGaapListenerStatisticsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGaapListenerStatisticsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_gaap_listener_statistics.listener_statistics"),
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
