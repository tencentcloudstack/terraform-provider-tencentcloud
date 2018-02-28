package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceTencentCloudSubnet_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: TestAccDataSourceTencentCloudSubnetConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_subnet.foo"),
					resource.TestCheckResourceAttr("data.tencentcloud_subnet.foo", "name", "ci-terraform-subnet-do-not-delete"),
					resource.TestCheckResourceAttr("data.tencentcloud_subnet.foo", "availability_zone", "ap-guangzhou-3"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudSubnetConfig = `
data "tencentcloud_subnet" "foo" {
	vpc_id = "vpc-8ek64x3d"
	subnet_id = "subnet-b1wk8b10"
}
`
