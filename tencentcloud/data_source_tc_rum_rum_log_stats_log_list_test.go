package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudRumRumLogStatsLogListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumRumLogStatsLogListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_rum_log_stats_log_list.rum_log_stats_log_list")),
			},
		},
	})
}

const testAccRumRumLogStatsLogListDataSource = `

data "tencentcloud_rum_rum_log_stats_log_list" "rum_log_stats_log_list" {
  start_time = 1625444040
  query = "id:123 AND type:&quot;log&quot;"
  end_time = 1625454840
  i_d = 1
  }

`
