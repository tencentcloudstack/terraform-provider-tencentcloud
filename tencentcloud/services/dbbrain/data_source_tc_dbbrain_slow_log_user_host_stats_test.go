package dbbrain_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainSlowLogUserHostStatsDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -1).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainSlowLogUserHostStatsDataSource, tcacctest.DefaultDbBrainInstanceId, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_slow_log_user_host_stats.test"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_log_user_host_stats.test", "start_time", startTime),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_log_user_host_stats.test", "end_time", endTime),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_log_user_host_stats.test", "product", "mysql"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_user_host_stats.test", "items.#"),
				),
			},
		},
	})
}

const testAccDbbrainSlowLogUserHostStatsDataSource = `

data "tencentcloud_dbbrain_slow_log_user_host_stats" "test" {
  instance_id = "%s"
  start_time = "%s"
  end_time = "%s"
  product = "mysql"
}

`
