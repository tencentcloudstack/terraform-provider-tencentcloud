package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudTestingVpcV3Update(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTestingVpcConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("tencentcloud_vpc.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", defaultVpcCidr),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", defaultInsName),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "is_multicast", "true"),

					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "dns_servers.#"),
				),
			},
			{
				Config: testAccTestingVpcConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckVpcExists("tencentcloud_vpc.foo"),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "cidr_block", defaultVpcCidrLess),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "name", defaultInsNameUpdate),
					resource.TestCheckResourceAttr("tencentcloud_vpc.foo", "is_multicast", "false"),

					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "is_default"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "create_time"),
					resource.TestCheckResourceAttrSet("tencentcloud_vpc.foo", "dns_servers.#"),
				),
			},
		},
	})
}

const testAccTestingVpcConfig = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name
  cidr_block = var.vpc_cidr
}
`

const testAccTestingVpcConfigUpdate = defaultVpcVariable + `
resource "tencentcloud_vpc" "foo" {
  name       = var.instance_name_update
  cidr_block = var.vpc_cidr_less

  is_multicast = false
}
`
