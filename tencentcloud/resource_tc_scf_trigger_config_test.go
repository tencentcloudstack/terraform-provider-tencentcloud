package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudScfTriggerConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfTriggerConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_scf_trigger_config.trigger_config", "id")),
			},
			{
				ResourceName:      "tencentcloud_scf_trigger_config.trigger_config",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccScfTriggerConfig = `

resource "tencentcloud_scf_trigger_config" "trigger_config" {
  enable        = "OPEN"
  function_name = "keep-1676351130"
  trigger_name  = "SCF-timer-1685540160"
  type          = "timer"
  qualifier     = "$DEFAULT"
  namespace     = "default"
}

`
