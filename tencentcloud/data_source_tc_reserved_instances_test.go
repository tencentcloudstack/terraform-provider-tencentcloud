package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudReservedInstancesDataSource(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckCommon(t, ACCOUNT_TYPE_INTERNATION) },
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
  availability_zone = "ap-guangzhou-3"
}
`
