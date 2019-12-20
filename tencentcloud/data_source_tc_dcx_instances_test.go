package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceTencentCloudDcxV3InstancesBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudDcxInstances,
				Check: resource.ComposeTestCheckFunc(
					//name filter
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcx_instances.name_select"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcx_instances.name_select", "instance_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcx_instances.name_select", "name"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudDcxInstances = `
data "tencentcloud_dcx_instances"  "name_select" {
    name ="a"
}
`
