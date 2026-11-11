package captcha_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudCaptchaInfoInternationalResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCaptchaInfoInternational,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_captcha_info_international.example", "id"),
					resource.TestCheckResourceAttr("tencentcloud_captcha_info_international.example", "app_name", "tf-example"),
					resource.TestCheckResourceAttr("tencentcloud_captcha_info_international.example", "channel_info", "web"),
				),
			},
			{
				Config: testAccCaptchaInfoInternationalUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("tencentcloud_captcha_info_international.example", "app_name", "tf-example-update"),
				),
			},
			{
				ResourceName:      "tencentcloud_captcha_info_international.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCaptchaInfoInternational = `
resource "tencentcloud_captcha_info_international" "example" {
  app_name     = "tf-example"
  channel_info = "web"
  verify_rank  = "1"
}
`

const testAccCaptchaInfoInternationalUpdate = `
resource "tencentcloud_captcha_info_international" "example" {
  app_name     = "tf-example-update"
  channel_info = "web"
  verify_rank  = "2"
}
`
