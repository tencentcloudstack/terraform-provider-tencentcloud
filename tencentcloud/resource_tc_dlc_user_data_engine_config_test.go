package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcUserDataEngineConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcUserDataEngineConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_dlc_user_data_engine_config.user_data_engine_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDlcUserDataEngineConfig = `

resource "tencentcloud_dlc_user_data_engine_config" "user_data_engine_config" {
  data_engine_id = "DataEngine-g5ds87d8"
  data_engine_config_pairs {
		config_item = "key"
		config_value = "value"

  }
  session_resource_template {
		driver_size = "small"
		executor_size = "small"
		executor_nums = 1
		executor_max_numbers = 1

  }
}

`
