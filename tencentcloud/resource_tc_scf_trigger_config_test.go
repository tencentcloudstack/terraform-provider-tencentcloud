package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
  enable = ""
  function_name = ""
  trigger_name = ""
  type = ""
  qualifier = ""
  namespace = ""
  trigger_desc = ""
}

`
