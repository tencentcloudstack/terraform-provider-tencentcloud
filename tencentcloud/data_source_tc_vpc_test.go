package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudVpc_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudVpcConfig_id,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc.id"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpc.id", "name", "guagua_vpc_instance_test"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudVpcConfig_id = `
resource "tencentcloud_vpc" "foo" {
    name       = "guagua_vpc_instance_test"
    cidr_block = "10.0.0.0/16"
}

data "tencentcloud_vpc" "id" {
	id = "${tencentcloud_vpc.foo.id}"
}
`

