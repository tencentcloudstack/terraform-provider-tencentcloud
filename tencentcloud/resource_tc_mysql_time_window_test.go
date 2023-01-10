package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCdbTimeWindowResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCdbTimeWindow,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mysql_time_window.time_window", "id")),
			},
			{
				ResourceName:      "tencentcloud_mysql_time_window.time_window",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbTimeWindow = `

resource "tencentcloud_mysql_time_window" "time_window" {
  instance_id    = "cdb-fitq5t9h"
  max_delay_time = 10
  time_ranges    = [
    "01:00-02:01"
  ]
  weekdays       = [
    "friday",
    "monday",
    "saturday",
    "thursday",
    "tuesday",
    "wednesday",
  ]
}

`
