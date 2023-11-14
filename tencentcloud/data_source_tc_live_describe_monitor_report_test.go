package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveDescribeMonitorReportDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveDescribeMonitorReportDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_live_describe_monitor_report.describe_monitor_report")),
			},
		},
	})
}

const testAccLiveDescribeMonitorReportDataSource = `

data "tencentcloud_live_describe_monitor_report" "describe_monitor_report" {
  monitor_id = ""
    }

`
