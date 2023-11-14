package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsScheduleConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsScheduleConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_schedule_config.schedule_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_schedule_config.schedule_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsScheduleConfig = `

resource "tencentcloud_mps_schedule_config" "schedule_config" {
  schedule_id = 
}

`
