package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudInstanceTypesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceTypesDataSourceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.t4c8g", "instance_types.0.cpu_core_count", "4"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.t4c8g", "instance_types.0.memory_size", "8"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.t4c8g", "instance_types.0.availability_zone", "ap-guangzhou-3"),
				),
			},
		},
	})
}

func TestAccTencentCloudInstanceTypesDataSource_sell(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTencentCloudInstanceTypesDataSourceConfigSell,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.t4c8g", "instance_types.0.cpu_core_count", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.t4c8g", "instance_types.0.memory_size", "1"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.t4c8g", "instance_types.0.availability_zone", "ap-guangzhou-3"),
					resource.TestCheckResourceAttr("data.tencentcloud_instance_types.t4c8g", "instance_types.0.family", "SA2"),
				),
			},
		},
	})
}

const testAccTencentCloudInstanceTypesDataSourceConfigBasic = `
data "tencentcloud_instance_types" "t4c8g" {
  availability_zone = "ap-guangzhou-3"
  cpu_core_count = 4
  memory_size    = 8
}
`

const testAccTencentCloudInstanceTypesDataSourceConfigSell = `
data "tencentcloud_instance_types" "t4c8g" {
  cpu_core_count = 1
  memory_size    = 1
  exclude_sold_out = true

  filter{
	name = "instance-family"
    values = ["SA2"]
  }

  filter{
	name = "zone"
    values = ["ap-guangzhou-3"]
  }
}
`
