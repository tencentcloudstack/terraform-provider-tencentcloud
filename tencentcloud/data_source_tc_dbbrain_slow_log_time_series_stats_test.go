package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainSlowLogTimeSeriesStatsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainSlowLogTimeSeriesStatsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_slow_log_time_series_stats.slow_log_time_series_stats")),
			},
		},
	})
}

const testAccDbbrainSlowLogTimeSeriesStatsDataSource = `

data "tencentcloud_dbbrain_slow_log_time_series_stats" "slow_log_time_series_stats" {
  instance_id = ""
  start_time = ""
  end_time = ""
  product = ""
      }

`
