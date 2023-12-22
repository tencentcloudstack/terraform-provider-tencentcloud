package rum_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumReportCountDataSource_basic -v
func TestAccTencentCloudRumReportCountDataSource_basic(t *testing.T) {
	t.Parallel()

	startTime := time.Now().AddDate(0, 0, -29).Unix()
	endTime := time.Now().Unix()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccRumReportCountDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_rum_report_count.report_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_report_count.report_count", "result"),
				),
			},
		},
	})
}

const testAccRumReportCountDataSource = `

data "tencentcloud_rum_report_count" "report_count" {
  start_time  = %v
  end_time    = %v
  project_id  = 120000
  report_type = "log"
}

`
