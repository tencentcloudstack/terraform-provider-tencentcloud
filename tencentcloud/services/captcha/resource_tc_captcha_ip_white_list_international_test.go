package captcha_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCaptchaIpWhiteListInternationalResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCaptchaIpWhiteListInternational,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_captcha_ip_white_list_international.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_captcha_ip_white_list_international.example", "name", "tf-example"),
				),
			},
			{
				Config: testAccCaptchaIpWhiteListInternationalUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_captcha_ip_white_list_international.example", "comment", "updated"),
				),
			},
			{
				ResourceName:      "tencentcloud_captcha_ip_white_list_international.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCaptchaIpWhiteListInternational = `
resource "tencentcloud_captcha_ip_white_list_international" "example" {
  name          = "tf-example"
  captcha_appid = 179000003
  ip            = "1.2.3.4"
  comment       = "tf example"
}
`

const testAccCaptchaIpWhiteListInternationalUpdate = `
resource "tencentcloud_captcha_ip_white_list_international" "example" {
  name          = "tf-example"
  captcha_appid = 179000003
  ip            = "1.2.3.4"
  comment       = "updated"
}
`
