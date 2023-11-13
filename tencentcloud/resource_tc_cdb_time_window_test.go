package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cdb_time_window.time_window", "id")),
			},
			{
				ResourceName:      "tencentcloud_cdb_time_window.time_window",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCdbTimeWindow = `

resource "tencentcloud_cdb_time_window" "time_window" {
  instance_id = "cdb-c1nl9rpv"
  time_ranges = 
  weekdays = 
  max_delay_time = 
}

`
