package dayuv2_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosOverviewCcTrendDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosOverviewCcTrendDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_antiddos_overview_cc_trend.overview_cc_trend"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_overview_cc_trend.overview_cc_trend", "data.#"),
				),
			},
		},
	})
}

const testAccAntiddosOverviewCcTrendDataSource = `

data "tencentcloud_antiddos_overview_cc_trend" "overview_cc_trend" {
	period = 300
	start_time = "2023-11-20 00:00:00"
	end_time = "2023-11-21 00:00:00"
	metric_name = "inqps"
	business = "bgpip"
}
`
