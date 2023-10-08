package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudMpsEnableScheduleConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsEnableScheduleConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_enable_schedule_config.enable_schedule_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_enable_schedule_config.enable_schedule_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsEnableScheduleConfig = `

resource "tencentcloud_mps_enable_schedule_config" "enable_schedule_config" {
  schedule_id = 
}

`
