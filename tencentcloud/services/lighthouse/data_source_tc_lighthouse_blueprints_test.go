package lighthouse_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseBlueprintsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseBlueprintsDataSource,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_blueprints.blueprints"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_lighthouse_blueprints.blueprints", "blueprint_set.#"),
				),
			},
		},
	})
}

const testAccLighthouseBlueprintsDataSource = `
data "tencentcloud_lighthouse_blueprints" "blueprints" {
}
`

func TestAccTencentCloudLighthouseBlueprintsDataSource_filter(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseBlueprintsDataSourceFilter,
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_blueprints.linux"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_lighthouse_blueprints.linux", "blueprint_set.#"),
				),
			},
		},
	})
}

const testAccLighthouseBlueprintsDataSourceFilter = `
data "tencentcloud_lighthouse_blueprints" "linux" {
  filters {
    name   = "platform-type"
    values = ["LINUX_UNIX"]
  }
}
`
