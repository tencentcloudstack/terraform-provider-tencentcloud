package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLighthouseModifyInstanceRenewFlagResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLighthouseModifyInstanceRenewFlag,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_lighthouse_modify_instance_renew_flag.modify_instance_renew_flag", "id")),
			},
			{
				ResourceName:      "tencentcloud_lighthouse_modify_instance_renew_flag.modify_instance_renew_flag",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccLighthouseModifyInstanceRenewFlag = `

resource "tencentcloud_lighthouse_modify_instance_renew_flag" "modify_instance_renew_flag" {
  instance_ids = 
  renew_flag = "NOTIFY_AND_MANUAL_RENEW"
}

`
