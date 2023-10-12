package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumStaticUrlDataSource_basic -v
func TestAccTencentCloudRumStaticUrlDataSource_basic(t *testing.T) {
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
				Config: fmt.Sprintf(testAccRumStaticUrlDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_static_url.static_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_static_url.static_url", "result"),
				),
			},
		},
	})
}

const testAccRumStaticUrlDataSource = `

data "tencentcloud_rum_static_url" "static_url" {
  start_time = %v
  type       = "pagepv"
  end_time   = %v
  project_id = 120000
}

`
