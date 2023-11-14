package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseResetInstanceBlueprintDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseResetInstanceBlueprintDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_reset_instance_blueprint.reset_instance_blueprint")),
			},
		},
	})
}

const testAccLighthouseResetInstanceBlueprintDataSource = `

data "tencentcloud_lighthouse_reset_instance_blueprint" "reset_instance_blueprint" {
  instance_id = "lhins-123456"
  offset = 0
  limit = 20
  filters {
		name = "blueprint-id"
		values = 

  }
}

`
