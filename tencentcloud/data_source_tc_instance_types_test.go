package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudInstanceTypesDataSource_basic(t *testing.T) {
	var currentRegion = "ap-guangzhou"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreSetRegion(currentRegion)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceTypesDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_instance_types.i4c8g"),

					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.i4c8g", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.i4c8g", "instance_types.0.memory_size", "8"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.i4c8g", "instance_types.0.family", "S1"),
				),
			},
		},
	})
}

const testAccTencentCloudInstanceTypesDataSourceConfigBasic = `
data "tencentcloud_instance_types" "i4c8g" {
  filter {
    name = "zone"
    values = ["ap-guangzhou-3"]
  }
  filter {
    name = "instance-family"
    values = ["S1"]
  }
  cpu_core_count = 4
  memory_size = 8
}
`
