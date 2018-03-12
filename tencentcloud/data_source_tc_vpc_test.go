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
					resource.TestCheckResourceAttr("data.tencentcloud_vpc.id", "name", "ci-terraform-vpc-do-not-delete"),
				),
			},
			{
				Config: TestAccDataSourceTencentCloudVpcConfig_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_vpc.name"),
					resource.TestCheckResourceAttr("data.tencentcloud_vpc.name", "id", "vpc-8ek64x3d"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudVpcConfig_id = `
data "tencentcloud_vpc" "id" {
	id = "vpc-8ek64x3d"
}
`

const TestAccDataSourceTencentCloudVpcConfig_name = `
data "tencentcloud_vpc" "name" {
    name = "ci-terraform"
}
`
