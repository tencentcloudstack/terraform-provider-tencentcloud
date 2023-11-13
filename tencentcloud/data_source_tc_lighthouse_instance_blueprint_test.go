package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseInstanceBlueprintDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
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
  instance_ids = 
  }

`
