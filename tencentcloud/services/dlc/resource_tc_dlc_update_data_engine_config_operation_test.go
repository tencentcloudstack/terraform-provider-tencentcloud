package dlc_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDlcUpdateDataEngineConfigOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcUpdateDataEngineConfigOperation,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_update_data_engine_config_operation.update_data_engine_config_operation", "id"),
					resource.TestCheckResourceAttr("tencentcloud_dlc_update_data_engine_config_operation.update_data_engine_config_operation", "data_engine_id", "DataEngine-3t1r4f20"),
				),
			},
		},
	})
}

const testAccDlcUpdateDataEngineConfigOperation = `

resource "tencentcloud_dlc_update_data_engine_config_operation" "update_data_engine_config_operation" {
  data_engine_id = "DataEngine-3t1r4f20"
  data_engine_config_command = "UpdateSparkSQLLakefsPath"
}

`
