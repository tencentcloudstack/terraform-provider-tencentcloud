package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseInstanceBlueprintDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { tcacctest.AccPreCheckCommon(t, tcacctest.ACCOUNT_TYPE_PREPAY) },
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseInstanceBlueprintDataSource,
				Check:  resource.ComposeTestCheckFunc(tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_instance_blueprint.instance_blueprint")),
			},
		},
	})
}

const testAccLighthouseInstanceBlueprintDataSource = tcacctest.DefaultLighthoustVariables + `

data "tencentcloud_lighthouse_instance_blueprint" "instance_blueprint" {
	instance_ids = [var.lighthouse_id]
}
`
