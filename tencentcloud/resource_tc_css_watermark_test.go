package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCssWatermark_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssWatermark,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_watermark.watermark",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssWatermark = `

resource "tencentcloud_css_watermark" "watermark" {
  picture_url = ""
  watermark_name = ""
  x_position = ""
  y_position = ""
  width = ""
  height = ""
}

`
