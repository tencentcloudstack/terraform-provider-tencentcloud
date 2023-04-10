package tencentcloud

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudLighthouseBlueprintResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseBlueprint,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_blueprint.blueprint", "id")),
			},
			{
				PreConfig:               func() { time.Sleep(time.Minute * 1) },
				ResourceName:            "tencentcloud_lighthouse_blueprint.blueprint",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"instance_id"},
			},
		},
	})
}

const testAccLighthouseBlueprint = `

resource "tencentcloud_lighthouse_blueprint" "blueprint" {
  blueprint_name = "blueprint_name_test"
  description = "blueprint_description_test"
  instance_id = "lhins-hwe21u91"
}

`
