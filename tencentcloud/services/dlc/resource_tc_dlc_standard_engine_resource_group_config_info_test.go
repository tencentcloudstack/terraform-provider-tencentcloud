package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcStandardEngineResourceGroupConfigInfoResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcStandardEngineResourceGroupConfigInfo,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group_config_info.example", "id"),
				),
			},
			{
				Config: testAccDlcStandardEngineResourceGroupConfigInfoUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group_config_info.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_dlc_standard_engine_resource_group_config_info.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDlcStandardEngineResourceGroupConfigInfo = `
resource "tencentcloud_dlc_standard_engine_resource_group_config_info" "example" {
  engine_resource_group_name = "tf-example"
  static_conf_context {
    params {
      config_item  = "item1"
      config_value = "value1"
    }

    params {
      config_item  = "item2"
      config_value = "value2"
    }
  }

  dynamic_conf_context {
    params {
      config_item  = "item3"
      config_value = "value3"
    }
  }
}

`

const testAccDlcStandardEngineResourceGroupConfigInfoUpdate = `
resource "tencentcloud_dlc_standard_engine_resource_group_config_info" "example" {
  engine_resource_group_name = "tf-example"
  static_conf_context {
    params {
      config_item  = "item1"
      config_value = "value1"
    }
  }

  dynamic_conf_context {
    params {
      config_item  = "item3"
      config_value = "value3"
    }
	
	params {
      config_item  = "item2"
      config_value = "value2"
    }
  }
}
`
