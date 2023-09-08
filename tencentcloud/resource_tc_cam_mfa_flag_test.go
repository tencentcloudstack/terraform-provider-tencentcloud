package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
data "tencentcloud_user_info" "info"{}

resource "tencentcloud_cam_mfa_flag" "mfa_flag" {
  op_uin = data.tencentcloud_user_info.info.uin
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
