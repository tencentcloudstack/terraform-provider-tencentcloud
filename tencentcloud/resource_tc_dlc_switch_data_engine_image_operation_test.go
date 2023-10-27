package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDlcSwitchDataEngineImageOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcSwitchDataEngineImageOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_switch_data_engine_image_operation.switch_data_engine_image_operation", "id")),
			},
		},
	})
}

const testAccDlcSwitchDataEngineImageOperation = `

resource "tencentcloud_dlc_switch_data_engine_image_operation" "switch_data_engine_image_operation" {
  data_engine_id = "DataEngine-cgkvbas6"
  new_image_version_id = "f54fba71-5f9c-4dfe-a565-004d7b6d3864"
}

`
