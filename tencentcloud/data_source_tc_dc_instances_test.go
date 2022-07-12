package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceTencentCloudDcV3InstancesBasic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudDcInstances,
				Check: resource.ComposeTestCheckFunc(
					//name filter
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dc_instances.name_select"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dc_instances.name_select", "instance_list.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dc_instances.name_select", "name"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudDcInstances = `
data tencentcloud_dc_instances  name_select {
    name ="x"
}
`
