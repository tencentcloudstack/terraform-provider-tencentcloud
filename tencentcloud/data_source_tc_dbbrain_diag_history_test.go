package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudDbbrainDiagHistoryDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -1).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainDiagHistoryDataSource, defaultDbBrainInstanceId, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_diag_history.diag_history"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_history.diag_history", "events.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_history.diag_history", "events.0.diag_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_history.diag_history", "events.0.event_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_history.diag_history", "events.0.start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_history.diag_history", "events.0.end_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_history.diag_history", "events.0.event_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_history.diag_history", "events.0.severity"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_history.diag_history", "events.0.outline"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_history.diag_history", "events.0.diag_item"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_history.diag_history", "events.0.instance_id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_history.diag_history", "events.0.region"),
				),
			},
		},
	})
}

const testAccDbbrainDiagHistoryDataSource = `

data "tencentcloud_dbbrain_diag_history" "diag_history" {
  instance_id = "%s"
  start_time = "%s"
  end_time = "%s"
  product = "mysql"
}

`
