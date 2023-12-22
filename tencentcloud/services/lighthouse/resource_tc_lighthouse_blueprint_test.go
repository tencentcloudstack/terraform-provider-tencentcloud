package lighthouse_test

import (
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseBlueprintResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
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

const testAccLighthouseBlueprint = tcacctest.DefaultLighthoustVariables + `

resource "tencentcloud_lighthouse_blueprint" "blueprint" {
  blueprint_name = "blueprint_name_test"
  description = "blueprint_description_test"
  instance_id = var.lighthouse_id
}

`
