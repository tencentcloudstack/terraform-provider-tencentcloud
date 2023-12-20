package cam_test

import (
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCamMfaFlagResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCamMfaFlag,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cam_mfa_flag.mfa_flag", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_mfa_flag.mfa_flag", "login_flag.#"),
					resource.TestCheckResourceAttr("tencentcloud_cam_mfa_flag.mfa_flag", "login_flag.0.phone", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cam_mfa_flag.mfa_flag", "login_flag.0.stoken", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cam_mfa_flag.mfa_flag", "login_flag.0.wechat", "0"),

					resource.TestCheckResourceAttrSet("tencentcloud_cam_mfa_flag.mfa_flag", "action_flag.#"),
					resource.TestCheckResourceAttr("tencentcloud_cam_mfa_flag.mfa_flag", "action_flag.0.phone", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cam_mfa_flag.mfa_flag", "action_flag.0.stoken", "1"),
					resource.TestCheckResourceAttr("tencentcloud_cam_mfa_flag.mfa_flag", "action_flag.0.wechat", "0"),
				),
			},
			{
				Config: testAccCamMfaFlagUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cam_mfa_flag.mfa_flag", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cam_mfa_flag.mfa_flag", "login_flag.#"),
					resource.TestCheckResourceAttr("tencentcloud_cam_mfa_flag.mfa_flag", "login_flag.0.phone", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cam_mfa_flag.mfa_flag", "login_flag.0.stoken", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cam_mfa_flag.mfa_flag", "login_flag.0.wechat", "0"),

					resource.TestCheckResourceAttrSet("tencentcloud_cam_mfa_flag.mfa_flag", "action_flag.#"),
					resource.TestCheckResourceAttr("tencentcloud_cam_mfa_flag.mfa_flag", "action_flag.0.phone", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cam_mfa_flag.mfa_flag", "action_flag.0.stoken", "0"),
					resource.TestCheckResourceAttr("tencentcloud_cam_mfa_flag.mfa_flag", "action_flag.0.wechat", "0"),
				),
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
const testAccCamMfaFlagUpdate = `
data "tencentcloud_user_info" "info"{}

resource "tencentcloud_cam_mfa_flag" "mfa_flag" {
  op_uin = data.tencentcloud_user_info.info.uin
  login_flag {
	phone = 0
	stoken = 0
	wechat = 0
  }
  action_flag {
	phone = 0
	stoken = 0
	wechat = 0
  }
}

`
