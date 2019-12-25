package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
					resource.TestCheckResourceAttr("data.tencentcloud_vpc.id", "name", "tf-ci-test"),
				),
			},
		},
	})
}

const TestAccDataSourceTencentCloudVpcConfig_id = `
resource "tencentcloud_vpc" "foo" {
  name       = "tf-ci-test"
  cidr_block = "10.0.0.0/16"
}

data "tencentcloud_vpc" "id" {
  id = tencentcloud_vpc.foo.id
}
`
