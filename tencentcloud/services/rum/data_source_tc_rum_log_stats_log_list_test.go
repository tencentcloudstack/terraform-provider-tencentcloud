package rum_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudNeedFixRumLogStatsLogListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumLogStatsLogListDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_rum_log_stats_log_list.log_stats_log_list")),
			},
		},
	})
}

const testAccRumLogStatsLogListDataSource = `

data "tencentcloud_rum_log_stats_log_list" "log_stats_log_list" {
  start_time = 1625444040
  query = "id:123 AND type:\"log\""
  end_time = 1625454840
  project_id = 1
}

`
