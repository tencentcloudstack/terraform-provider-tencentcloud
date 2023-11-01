package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCssMonitorReportDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssMonitorReportDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_css_monitor_report.monitor_report"),
				),
			},
		},
	})
}

const testAccCssMonitorReportDataSource = `

data "tencentcloud_css_monitor_report" "monitor_report" {
	monitor_id = "0e8a12b5-df2a-4a1b-aa98-97d5610aa142"
}

`
