package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCvmChcDeployVpcResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcDeployVpc,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_deploy_vpc.chc_deploy_vpc", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_chc_deploy_vpc.chc_deploy_vpc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmChcDeployVpc = `
resource "tencentcloud_cvm_chc_deploy_vpc" "chc_deploy_vpc" {
	chc_id = "chc-0brmw3wl"
	deploy_virtual_private_cloud {
		  vpc_id = "vpc-4owdpnwr"
		  subnet_id = "subnet-j56j1u5u"
	}
	deploy_security_group_ids = ["sg-ijato2x1"]
}
`
