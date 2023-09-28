package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwNatInstanceResource_basic -v
func TestAccTencentCloudNeedFixCfwNatInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwNatInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_instance.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cfw_nat_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCfwNatInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_nat_instance.example", "id"),
				),
			},
		},
	})
}

const testAccCfwNatInstance = `
resource "tencentcloud_cfw_nat_instance" "example" {
  name  = "tf_example"
  width = 20
  mode  = 0
  new_mode_items {
    vpc_list = [
      "vpc-5063ta4i"
    ]
    eips = [
      "152.136.168.192"
    ]
  }
  cross_a_zone = 0
  zone_set     = [
    "ap-guangzhou-7"
  ]
}
`

const testAccCfwNatInstanceUpdate = `
resource "tencentcloud_cfw_nat_instance" "example" {
  name  = "tf_example_update"
  width = 20
  mode  = 0
  new_mode_items {
    vpc_list = [
      "vpc-5063ta4i"
    ]
    eips = [
      "152.136.168.192"
    ]
  }
  cross_a_zone = 0
  zone_set     = [
    "ap-guangzhou-7"
  ]
}
`
