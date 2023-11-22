package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosOverviewAttackTrendDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosOverviewAttackTrendDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_antiddos_overview_attack_trend.overview_attack_trend"),
					resource.TestCheckResourceAttr("data.tencentcloud_antiddos_overview_attack_trend.overview_attack_trend", "type", "ddos"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_attack_trend.overview_attack_trend", "start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_attack_trend.overview_attack_trend", "end_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_antiddos_overview_attack_trend.overview_attack_trend", "period", "86400"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_attack_trend.overview_attack_trend", "data.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_attack_trend.overview_attack_trend", "period_point_count"),
				),
			},
		},
	})
}

const testAccAntiddosOverviewAttackTrendDataSource = `

data "tencentcloud_antiddos_overview_attack_trend" "overview_attack_trend" {
  type = "ddos"
  dimension = "attackcount"
  period = 86400
  start_time = "2023-11-21 10:28:31"
  end_time = "2023-11-22 10:28:31"
}

`
