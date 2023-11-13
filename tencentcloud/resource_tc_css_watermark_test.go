package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCssWatermarkResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCssWatermark,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_css_watermark.watermark", "id")),
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
  picture_url = &lt;nil&gt;
  watermark_name = &lt;nil&gt;
  x_position = &lt;nil&gt;
  y_position = &lt;nil&gt;
  width = &lt;nil&gt;
  height = &lt;nil&gt;
}

`
