package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccTencentCloudReservedInstancesDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccReservedInstancesDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.tencentcloud_reserved_instances.instances", "reserved_instance_list.#"),
				),
			},
		},
	})
}

const testAccReservedInstancesDataSource = `
data "tencentcloud_reserved_instances" "instances" {
  availability_zone = "ap-guangzhou-2"
}
`
