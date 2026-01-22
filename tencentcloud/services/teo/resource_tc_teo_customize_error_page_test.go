package teo_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"
)

func TestAccTencentCloudTeoCustomizeErrorPageResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoCustomizeErrorPage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "content_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "content"),
				),
			},
			{
				Config: testAccTeoCustomizeErrorPageUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "zone_id"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "name"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "content_type"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "description"),
					resource.TestCheckResourceAttrSet("tencentcloud_teo_customize_error_page.example", "content"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_customize_error_page.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoCustomizeErrorPage = `
resource "tencentcloud_teo_customize_error_page" "example" {
  zone_id      = "zone-3edjdliiw3he"
  name         = "tf-example"
  content_type = "text/plain"
  description  = "description."
  content      = "customize error page"
}
`

const testAccTeoCustomizeErrorPageUpdate = `
resource "tencentcloud_teo_customize_error_page" "example" {
  zone_id      = "zone-3edjdliiw3he"
  name         = "tf-example-update"
  content_type = "application/json"
  description  = "description update."
  content = jsonencode({
    "key" : "value",
  })
}
`
