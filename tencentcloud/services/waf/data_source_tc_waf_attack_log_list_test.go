package waf_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudWafAttackLogListDataSource_basic -v
func TestAccTencentCloudWafAttackLogListDataSource_basic(t *testing.T) {
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
				Config: fmt.Sprintf(testAccWafAttackLogListDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_waf_attack_log_list.example"),
				),
			},
		},
	})
}

const testAccWafAttackLogListDataSource = `
data "tencentcloud_waf_attack_log_list" "example" {
  domain       = "all"
  start_time   = "%s"
  end_time     = "%s"
  query_string = "method:GET"
  sort         = "desc"
  query_count  = 10
  page         = 0
}
`
