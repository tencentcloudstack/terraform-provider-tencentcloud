package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcUserDataEngineConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcUserDataEngineConfig,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "data_engine_id", "DataEngine-cgkvbas6"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "data_engine_config_pairs.#"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "data_engine_config_pairs.0.config_item", "qq"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "data_engine_config_pairs.0.config_value", "aa"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "session_resource_template.#"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "session_resource_template.0.driver_size", "small"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "session_resource_template.0.executor_size", "small"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "session_resource_template.0.executor_nums", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "session_resource_template.0.executor_max_numbers", "1"),
				),
			},
			{
				Config: testAccDlcUserDataEngineConfigUpdate,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "data_engine_id", "DataEngine-cgkvbas6"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "data_engine_config_pairs.#"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "data_engine_config_pairs.0.config_item", "qq"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "data_engine_config_pairs.0.config_value", "ff"),
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "session_resource_template.#"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "session_resource_template.0.driver_size", "small"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "session_resource_template.0.executor_size", "small"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "session_resource_template.0.executor_nums", "1"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_user_data_engine_config.user_data_engine_config", "session_resource_template.0.executor_max_numbers", "1"),
				),
			},
		},
	})
}

const testAccDlcUserDataEngineConfig = `

resource "tencentcloud_dlc_user_data_engine_config" "user_data_engine_config" {
  data_engine_id = "DataEngine-cgkvbas6"
  data_engine_config_pairs {
		config_item = "qq"
		config_value = "aa"
  }
  session_resource_template {
		driver_size = "small"
		executor_size = "small"
		executor_nums = 1
		executor_max_numbers = 1
  }
}

`
const testAccDlcUserDataEngineConfigUpdate = `

resource "tencentcloud_dlc_user_data_engine_config" "user_data_engine_config" {
  data_engine_id = "DataEngine-cgkvbas6"
  data_engine_config_pairs {
		config_item = "qq"
		config_value = "ff"
  }
  session_resource_template {
		driver_size = "small"
		executor_size = "small"
		executor_nums = 1
		executor_max_numbers = 1
  }
}

`
