package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudReservedInstanceConfigsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccReservedInstanceConfigsDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_reserved_instance_configs.configs", "config_list.#"),
				),
			},
		},
	})
}

const testAccReservedInstanceConfigsDataSource = `
data "tencentcloud_reserved_instance_configs" "configs" {
  availability_zone = "ap-guangzhou-2"
}
`
