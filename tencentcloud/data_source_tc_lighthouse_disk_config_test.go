package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudLighthouseDiskConfigDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseDiskConfigDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_lighthouse_disk_config.disk_config")),
			},
		},
	})
}

const testAccLighthouseDiskConfigDataSource = `

data "tencentcloud_lighthouse_disk_config" "disk_config" {
  filters {
	name = "zone"
	values = ["ap-guangzhou-3"]
  }
}
`
