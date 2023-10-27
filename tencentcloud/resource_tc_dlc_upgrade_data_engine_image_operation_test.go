package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcUpgradeDataEngineImageOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcUpgradeDataEngineImageOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_upgrade_data_engine_image_operation.upgrade_data_engine_image_operation", "id")),
			},
		},
	})
}

const testAccDlcUpgradeDataEngineImageOperation = `

resource "tencentcloud_dlc_upgrade_data_engine_image_operation" "upgrade_data_engine_image_operation" {
  data_engine_id = "DataEngine-o6cs25y2"
}

`
