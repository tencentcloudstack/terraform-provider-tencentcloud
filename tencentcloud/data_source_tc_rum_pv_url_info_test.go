package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumPvUrlInfoDataSource_basic -v
func TestAccTencentCloudRumPvUrlInfoDataSource_basic(t *testing.T) {
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
				Config: fmt.Sprintf(testAccRumPvUrlInfoDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_pv_url_info.pv_url_info"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_pv_url_info.pv_url_info", "result"),
				),
			},
		},
	})
}

const testAccRumPvUrlInfoDataSource = `

data "tencentcloud_rum_pv_url_info" "pv_url_info" {
  start_time = %v
  type       = "pagepv"
  end_time   = %v
  project_id = 120000
}

`
