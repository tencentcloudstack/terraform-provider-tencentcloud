package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudTeoCustomErrorPage_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccTeoCustomErrorPage,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_teo_custom_error_page.custom_error_page", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_teo_custom_error_page.customErrorPage",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccTeoCustomErrorPage = `

resource "tencentcloud_teo_custom_error_page" "custom_error_page" {
  zone_id = ""
  entity = ""
    name = ""
  content = ""
}

`
