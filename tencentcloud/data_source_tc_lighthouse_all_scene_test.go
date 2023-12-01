package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseAllSceneDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseAllSceneDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_all_scene.scene")),
			},
		},
	})
}

const testAccLighthouseAllSceneDataSource = `

data "tencentcloud_lighthouse_all_scene" "scene" {
  offset = 0
  limit = 20
}

`
