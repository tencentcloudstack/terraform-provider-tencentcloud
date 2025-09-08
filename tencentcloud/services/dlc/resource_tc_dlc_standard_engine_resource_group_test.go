package dlc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudDlcStandardEngineResourceGroupResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccDlcStandardEngineResourceGroup,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_dlc_standard_engine_resource_group.dlc_standard_engine_resource_group", "id")),
		}, {
			ResourceName:      "tencentcloud_dlc_standard_engine_resource_group.dlc_standard_engine_resource_group",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccDlcStandardEngineResourceGroup = `

resource "tencentcloud_dlc_standard_engine_resource_group" "dlc_standard_engine_resource_group" {
  static_config_pairs = {
  }
  dynamic_config_pairs = {
  }
}
`
