package rum_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumPvUrlStatisticsDataSource_basic -v
func TestAccTencentCloudRumPvUrlStatisticsDataSource_basic(t *testing.T) {
	t.Parallel()

	startTime := time.Now().AddDate(0, 0, -29).Unix()
	endTime := time.Now().Unix()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccRumPvUrlStatisticsDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_rum_pv_url_statistics.pv_url_statistics"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_pv_url_statistics.pv_url_statistics", "result"),
				),
			},
		},
	})
}

const testAccRumPvUrlStatisticsDataSource = `

data "tencentcloud_rum_pv_url_statistics" "pv_url_statistics" {
  start_time = %v
  type       = "allcount"
  end_time   = %v
  project_id = 120000
}

`
