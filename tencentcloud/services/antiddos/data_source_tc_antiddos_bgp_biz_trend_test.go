package antiddos_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudAntiddosBgpBizTrendDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiddosBgpBizTrendDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_antiddos_bgp_biz_trend.bgp_biz_trend"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_bgp_biz_trend.bgp_biz_trend", "data_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_bgp_biz_trend.bgp_biz_trend", "total"),
					resource.TestCheckResourceAttr("data.tencentcloud_antiddos_bgp_biz_trend.bgp_biz_trend", "metric_name", "intraffic"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_antiddos_bgp_biz_trend.bgp_biz_trend", "max_data"),
				),
			},
		},
	})
}

const testAccAntiddosBgpBizTrendDataSource = `

data "tencentcloud_antiddos_bgp_biz_trend" "bgp_biz_trend" {
  business = "bgp-multip"
  start_time = "2023-11-22 09:25:00"
  end_time = "2023-11-22 10:25:00"
  metric_name = "intraffic"
  instance_id = "bgp-00000ry7"
  flag = 0
}

`
