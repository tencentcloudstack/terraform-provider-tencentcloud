package cvm_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCvmChcConfigResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "instance_name", "test"),
				),
			},
			{
				Config: testAccCvmChcConfig_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_config.chc_config", "id"),
					resource.TestCheckResourceAttr("tencentcloud_cvm_chc_config.chc_config", "instance_name", "test_update"),
				),
			},
			{
				ResourceName:            "tencentcloud_cvm_chc_config.chc_config",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bmc_user", "password"},
			},
		},
	})
}

const testAccCvmChcConfig = `
resource "tencentcloud_cvm_chc_config" "chc_config" {
	chc_id = "chc-0brmw3wl"
	instance_name = "test"
	bmc_user = "admin"
	password = "123"
	bmc_virtual_private_cloud {
	  vpc_id = "vpc-4owdpnwr"
	  subnet_id = "subnet-j56j1u5u"
	}
	bmc_security_group_ids = ["sg-ijato2x1"]
  
	deploy_virtual_private_cloud {
	  vpc_id = "vpc-4owdpnwr"
	  subnet_id = "subnet-j56j1u5u"
	}
	deploy_security_group_ids = ["sg-ijato2x1"]
  }
`

const testAccCvmChcConfig_update = `
resource "tencentcloud_cvm_chc_config" "chc_config" {
	chc_id = "chc-0brmw3wl"
	instance_name = "test_update"
	bmc_user = "admin"
	password = "123123"
	bmc_virtual_private_cloud {
	  vpc_id = "vpc-4owdpnwr"
	  subnet_id = "subnet-j56j1u5u"
	}
	bmc_security_group_ids = ["sg-ijato2x1"]
  
	deploy_virtual_private_cloud {
	  vpc_id = "vpc-4owdpnwr"
	  subnet_id = "subnet-j56j1u5u"
	}
	deploy_security_group_ids = ["sg-ijato2x1"]
  }
`
