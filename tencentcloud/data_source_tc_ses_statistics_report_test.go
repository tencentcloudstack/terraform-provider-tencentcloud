package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudSesStatisticsReportDataSource_basic -v
func TestAccTencentCloudSesStatisticsReportDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccStepSetRegion(t, "ap-hongkong")
			testAccPreCheckBusiness(t, ACCOUNT_TYPE_SES)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesStatisticsReportDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_ses_statistics_report.statistics_report"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "daily_volumes.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "daily_volumes.0.accepted_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "daily_volumes.0.bounce_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "daily_volumes.0.clicked_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "daily_volumes.0.delivered_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "daily_volumes.0.opened_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "daily_volumes.0.request_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "daily_volumes.0.send_date"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "daily_volumes.0.unsubscribe_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "overall_volume.0.accepted_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "overall_volume.0.bounce_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "overall_volume.0.delivered_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "overall_volume.0.opened_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "overall_volume.0.request_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "overall_volume.0.send_date"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_ses_statistics_report.statistics_report", "overall_volume.0.unsubscribe_count"),
				),
			},
		},
	})
}

const testAccSesStatisticsReportDataSource = `

data "tencentcloud_ses_statistics_report" "statistics_report" {
  start_date = "2020-10-01"
  end_date = "2023-09-05"
  domain = "iac-tf.cloud"
}

`
