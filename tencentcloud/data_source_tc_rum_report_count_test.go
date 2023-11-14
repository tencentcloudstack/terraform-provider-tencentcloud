package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudRumReportCountDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccRumReportCountDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_report_count.report_count")),
			},
		},
	})
}

const testAccRumReportCountDataSource = `

data "tencentcloud_rum_report_count" "report_count" {
  start_time = 1625444040
  end_time = 1625454840
  i_d = 1
  report_type = "log"
  instance_i_d = "rum-xxx"
  }

`
