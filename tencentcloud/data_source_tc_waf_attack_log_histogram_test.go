package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafAttackLogHistogramDataSource_basic -v
func TestAccTencentCloudWafAttackLogHistogramDataSource_basic(t *testing.T) {
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
				Config: fmt.Sprintf(testAccWafAttackLogHistogramDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_attack_log_histogram.example"),
				),
			},
		},
	})
}

const testAccWafAttackLogHistogramDataSource = `
data "tencentcloud_waf_attack_log_histogram" "example" {
  domain       = "all"
  start_time   = "%s"
  end_time     = "%s"
  query_string = "method:GET"
}
`
