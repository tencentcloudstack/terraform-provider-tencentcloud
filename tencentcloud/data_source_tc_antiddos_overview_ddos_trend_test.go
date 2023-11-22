package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosOverviewDdosTrendDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosOverviewDdosTrendDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_antiddos_overview_ddos_trend.overview_ddos_trend"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_ddos_trend.overview_ddos_trend", "data.#"),
				),
			},
		},
	})
}

const testAccAntiddosOverviewDdosTrendDataSource = `

data "tencentcloud_antiddos_overview_ddos_trend" "overview_ddos_trend" {
  period = 300
  start_time = "2023-11-20 14:16:23"
  end_time = "2023-11-21 14:16:23"
  metric_name = "bps"
  business = "bgpip"
}
`
