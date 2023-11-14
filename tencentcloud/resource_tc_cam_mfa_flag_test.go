package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCamMfaFlagResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamMfaFlag,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_cam_mfa_flag.mfa_flag", "id")),
			},
			{
				ResourceName:      "tencentcloud_cam_mfa_flag.mfa_flag",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCamMfaFlag = `

resource "tencentcloud_cam_mfa_flag" "mfa_flag" {
  op_uin = 20003xxxxxxx
  login_flag {
		phone = 0
		stoken = 1
		wechat = 0

  }
  action_flag {
		phone = 0
		stoken = 1
		wechat = 0

  }
}

`
