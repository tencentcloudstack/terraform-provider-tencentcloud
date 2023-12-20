package dbbrain_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainDiagEventDataSource_basic(t *testing.T) {
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
				Config: fmt.Sprintf(testAccDbbrainDiagEventDataSource, tcacctest.DefaultDbBrainInstanceId, startTime, endTime, tcacctest.DefaultDbBrainInstanceId),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_diag_event.diag_event"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "diag_item"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "diag_type"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "explanation"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "outline"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "problem"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "severity"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "start_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "suggestions"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_diag_event.diag_event", "end_time"),
				),
			},
		},
	})
}

const testAccDbbrainDiagEventDataSource = `

data "tencentcloud_dbbrain_diag_history" "diag_history" {
	instance_id = "%s"
	start_time = "%s"
	end_time = "%s"
	product = "mysql"
}

data "tencentcloud_dbbrain_diag_event" "diag_event" {
  instance_id = "%s"
  event_id = data.tencentcloud_dbbrain_diag_history.diag_history.events.0.event_id
  product = "mysql"
}

`
