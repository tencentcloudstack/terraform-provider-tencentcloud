package cwp_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCwpAutoOpenProversionConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{{
			Config: testAccCwpAutoOpenProversionConfig,
			Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cwp_auto_open_proversion_config.cwp_auto_open_proversion_config", "id")),
		}, {
			ResourceName:      "tencentcloud_cwp_auto_open_proversion_config.cwp_auto_open_proversion_config",
			ImportState:       true,
			ImportStateVerify: true,
		}},
	})
}

const testAccCwpAutoOpenProversionConfig = `

resource "tencentcloud_cwp_auto_open_proversion_config" "cwp_auto_open_proversion_config" {
}
`
