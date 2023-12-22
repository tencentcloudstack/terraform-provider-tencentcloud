package waf_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafPeakPointsDataSource_basic -v
func TestAccTencentCloudWafPeakPointsDataSource_basic(t *testing.T) {
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
				Config: fmt.Sprintf(testAccWafPeakPointsDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_waf_peak_points.example"),
				),
			},
		},
	})
}

const testAccWafPeakPointsDataSource = `
data "tencentcloud_waf_peak_points" "example" {
  from_time   = "%s"
  to_time     = "%s"
  domain      = "test.com"
  edition     = "clb-waf"
  instance_id = "waf_2kxtlbky00b2v1fn"
  metric_name = "access"
}
`
