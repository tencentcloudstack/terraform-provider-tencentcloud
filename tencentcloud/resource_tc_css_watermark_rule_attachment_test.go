package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCssWatermarkRule_basic(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssWatermarkRule,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("tencentcloud_css_watermark_rule.watermark_rule", "id"),
				),
			},
			{
				ResourceName:      "tencentcloud_css_watermark_rule.watermarkRule",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCssWatermarkRule = `

resource "tencentcloud_css_watermark_rule" "watermark_rule" {
  domain_name = ""
  app_name = ""
  stream_name = ""
  watermark_id = ""
    }

`
