package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCvmChcAttributeResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCvmChcAttribute,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cvm_chc_attribute.chc_attribute", "id")),
			},
			{
				ResourceName:      "tencentcloud_cvm_chc_attribute.chc_attribute",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCvmChcAttribute = `

resource "tencentcloud_cvm_chc_attribute" "chc_attribute" {
  chc_ids = 
  instance_name = ""
  device_type = ""
  bmc_user = ""
  password = ""
  bmc_security_group_ids = 
}

`
