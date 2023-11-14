package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
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
  chc_ids = 
  deploy_virtual_private_cloud {
		vpc_id = ""
		subnet_id = ""
		as_vpc_gateway = 
		private_ip_addresses = 
		ipv6_address_count = 

  }
  deploy_security_group_ids = 
}

`
