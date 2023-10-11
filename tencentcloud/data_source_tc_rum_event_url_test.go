package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -test.run TestAccTencentCloudRumEventUrlDataSource_basic -v
func TestAccTencentCloudRumEventUrlDataSource_basic(t *testing.T) {
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
				Config: fmt.Sprintf(testAccRumEventUrlDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_rum_event_url.event_url"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_rum_event_url.event_url", "result"),
				),
			},
		},
	})
}

const testAccRumEventUrlDataSource = `

data "tencentcloud_rum_event_url" "event_url" {
  start_time = %v
  type       = "allcount"
  end_time   = %v
  project_id = 120000
}

`
