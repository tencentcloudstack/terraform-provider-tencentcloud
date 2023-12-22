package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseResetInstanceBlueprintDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseResetInstanceBlueprintDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_reset_instance_blueprint.reset_instance_blueprint")),
			},
		},
	})
}

const testAccLighthouseResetInstanceBlueprintDataSource = tcacctest.DefaultLighthoustVariables + `

data "tencentcloud_lighthouse_reset_instance_blueprint" "reset_instance_blueprint" {
  instance_id = var.lighthouse_id
  offset = 0
  limit = 20
}
`
