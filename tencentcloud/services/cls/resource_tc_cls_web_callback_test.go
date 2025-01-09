package cls_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudClsWebCallbackResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccClsWebCallback,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_web_callback.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_web_callback.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_web_callback.example", "type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_web_callback.example", "webhook"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_web_callback.example", "method"),
				),
			},
			{
				Config: testAccClsWebCallbackUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_cls_web_callback.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_web_callback.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_web_callback.example", "type"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_web_callback.example", "webhook"),
					resource.TestCheckResourceAttrSet("tencentcloud_cls_web_callback.example", "method"),
				),
			},
			{
				ResourceName:      "tencentcloud_cls_web_callback.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccClsWebCallback = `
resource "tencentcloud_cls_web_callback" "example" {
  name    = "tf-example"
  type    = "Http"
  webhook = "https://demo.com"
  method  = "POST"
}
`

const testAccClsWebCallbackUpdate = `
resource "tencentcloud_cls_web_callback" "example" {
  name    = "tf-example-update"
  type    = "Http"
  webhook = "https://demo.update.com"
  method  = "PUT"
}
`
