package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmChcAssistVpcResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcAssistVpc,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_assist_vpc.chc_assist_vpc", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_chc_assist_vpc.chc_assist_vpc",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmChcAssistVpc = `

resource "tencentcloud_cvm_chc_assist_vpc" "chc_assist_vpc" {
  chc_ids = 
  bmc_virtual_private_cloud {
		vpc_id = ""
		subnet_id = ""
		as_vpc_gateway = 
		private_ip_addresses = 
		ipv6_address_count = 

  }
  bmc_security_group_ids = 
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
