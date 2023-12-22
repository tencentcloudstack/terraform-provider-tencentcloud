package scf_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudScfTriggerConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccScfTriggerConfig,
				Check: resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_scf_trigger_config.trigger_config", "id"),
					resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("tencentcloud_scf_trigger_config.trigger_config", "enable", "OPEN")),
					resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("tencentcloud_scf_trigger_config.trigger_config", "function_name", "keep-1676351130")),
					resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("tencentcloud_scf_trigger_config.trigger_config", "trigger_name", "SCF-timer-1685540160")),
					resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("tencentcloud_scf_trigger_config.trigger_config", "type", "timer")),
					resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("tencentcloud_scf_trigger_config.trigger_config", "qualifier", "$DEFAULT")),
					resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("tencentcloud_scf_trigger_config.trigger_config", "namespace", "default")),
					resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("tencentcloud_scf_trigger_config.trigger_config", "trigger_desc", "* 1 2 * * * *")),
					resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("tencentcloud_scf_trigger_config.trigger_config", "description", "func")),
					resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("tencentcloud_scf_trigger_config.trigger_config", "custom_argument", "Information"))),
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
  trigger_desc = "* 1 2 * * * *"
  description = "func"
  custom_argument = "Information"
}

`
