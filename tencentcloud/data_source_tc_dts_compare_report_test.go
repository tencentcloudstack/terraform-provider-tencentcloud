package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDtsCompareReportDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDtsCompareReportDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dts_compare_report.compare_report")),
			},
		},
	})
}

const testAccDtsCompareReportDataSource = `

data "tencentcloud_dts_compare_report" "compare_report" {
  job_id = "dts-amm1jw5q"
  compare_task_id = "dts-amm1jw5q-cmp-bmuum7jk"
  difference_limit = 10
  difference_offset = 0
  difference_d_b = "db1"
  difference_table = "t1"
  skipped_limit = 10
  skipped_offset = 0
  skipped_d_b = "db1"
  skipped_table = "t1"
    }

`
