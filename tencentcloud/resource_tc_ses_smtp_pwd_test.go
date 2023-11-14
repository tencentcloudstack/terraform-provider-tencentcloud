package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudSesSmtp_pwdResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSesSmtp_pwd,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ses_smtp_pwd.smtp_pwd", "id")),
			},
			{
				ResourceName:      "tencentcloud_ses_smtp_pwd.smtp_pwd",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccSesSmtp_pwd = `

resource "tencentcloud_ses_smtp_pwd" "smtp_pwd" {
  password = "xX1@#xXXXX"
  email_address = "abc@ef.com"
}

`
