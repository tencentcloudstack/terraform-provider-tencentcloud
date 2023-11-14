package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainSlowLogUserHostStatsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainSlowLogUserHostStatsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_slow_log_user_host_stats.slow_log_user_host_stats")),
			},
		},
	})
}

const testAccDbbrainSlowLogUserHostStatsDataSource = `

data "tencentcloud_dbbrain_slow_log_user_host_stats" "slow_log_user_host_stats" {
  instance_id = ""
  start_time = ""
  end_time = ""
  product = ""
  md5 = ""
    }

`
