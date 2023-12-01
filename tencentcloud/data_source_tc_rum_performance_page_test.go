package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumPerformancePageDataSource_basic -v
func TestAccTencentCloudRumPerformancePageDataSource_basic(t *testing.T) {
	t.Parallel()

	startTime := time.Now().AddDate(0, 0, -29).Unix()
	endTime := time.Now().Unix()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccRumPerformancePageDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_performance_page.performance_page"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_performance_page.performance_page", "result"),
				),
			},
		},
	})
}

const testAccRumPerformancePageDataSource = `

data "tencentcloud_rum_performance_page" "performance_page" {
  start_time = %v
  type       = "pagepv"
  end_time   = %v
  project_id = 120000
}

`
