package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDlcRestartDataEngineOperationResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDlcRestartDataEngine,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_restart_data_engine_operation.restart_data_engine", "id")),
			},
		},
	})
}

const testAccDlcRestartDataEngine = `

resource "tencentcloud_dlc_restart_data_engine_operation" "restart_data_engine" {
  data_engine_id = "DataEngine-cgkvbas6"
  forced_operation = false
}

`
