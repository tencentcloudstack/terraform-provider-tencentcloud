package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesStatisticsReportDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesStatisticsReportDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_ses_statistics_report.statistics_report")),
			},
		},
	})
}

const testAccSesStatisticsReportDataSource = `

data "tencentcloud_ses_statistics_report" "statistics_report" {
  start_date = "2020-10-01"
  end_date = "2020-10-03"
  domain = "qcloud.com"
  receiving_mailbox_type = "gmail.com"
    }

`
