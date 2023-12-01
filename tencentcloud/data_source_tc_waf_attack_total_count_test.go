package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafAttackTotalCountDataSource_basic -v
func TestAccTencentCloudWafAttackTotalCountDataSource_basic(t *testing.T) {
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
				Config: fmt.Sprintf(testAccWafAttackTotalCountDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_waf_attack_total_count.example"),
				),
			},
		},
	})
}

const testAccWafAttackTotalCountDataSource = `
data "tencentcloud_waf_attack_total_count" "example" {
  start_time   = "%s"
  end_time     = "%s"
  domain       = "all"
  query_string = "method:GET"
}
`
