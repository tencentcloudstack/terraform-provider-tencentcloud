package cfw_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudNeedFixCfwVpcInstanceResource_basic -v
func TestAccTencentCloudNeedFixCfwVpcInstanceResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCfwVpcInstance,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_instance.example", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_cfw_vpc_instance.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccCfwVpcInstanceUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cfw_vpc_instance.example", "id"),
				),
			},
		},
	})
}

const testAccCfwVpcInstance = `
resource "tencentcloud_cfw_vpc_instance" "example" {
  name = "tf_example"
  mode = 0

  vpc_fw_instances {
    name    = "fw_ins_example"
    vpc_ids = [
      "vpc-291vnoeu",
      "vpc-39ixq9ci"
    ]
    fw_deploy {
      deploy_region = "ap-guangzhou"
      width         = 1024
      cross_a_zone  = 1
      zone_set      = [
        "ap-guangzhou-6",
        "ap-guangzhou-7"
      ]
    }
  }

  switch_mode = 1
  fw_vpc_cidr = "auto"
}
`

const testAccCfwVpcInstanceUpdate = `
resource "tencentcloud_cfw_vpc_instance" "example" {
  name = "tf_example_update"
  mode = 0

  vpc_fw_instances {
    name    = "fw_ins_example"
    vpc_ids = [
      "vpc-291vnoeu",
      "vpc-39ixq9ci"
    ]
    fw_deploy {
      deploy_region = "ap-guangzhou"
      width         = 1024
      cross_a_zone  = 1
      zone_set      = [
        "ap-guangzhou-6",
        "ap-guangzhou-7"
      ]
    }
  }

  switch_mode = 1
  fw_vpc_cidr = "auto"
}
`
