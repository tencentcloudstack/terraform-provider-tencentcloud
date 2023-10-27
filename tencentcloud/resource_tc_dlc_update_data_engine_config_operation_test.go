package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDlcUpdateDataEngineConfigOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcUpdateDataEngineConfigOperation,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_data_engine_config_operation.update_data_engine_config_operation", "id")),
			},
			{
				ResourceName:      "tencentcloud_dlc_update_data_engine_config_operation.update_data_engine_config_operation",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccDlcUpdateDataEngineConfigOperation = `

resource "tencentcloud_dlc_update_data_engine_config_operation" "update_data_engine_config_operation" {
  data_engine_ids = 
  data_engine_config_command = "UpdateSparkSQLLakefsPath"
}

`
