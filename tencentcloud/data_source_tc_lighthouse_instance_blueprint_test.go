package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseInstanceBlueprintDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_PREPAY) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseInstanceBlueprintDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_instance_blueprint.instance_blueprint")),
			},
		},
	})
}

const testAccLighthouseInstanceBlueprintDataSource = `

data "tencentcloud_lighthouse_instance_blueprint" "instance_blueprint" {
	instance_ids = ["lhins-hwe21u91"]
}
`
