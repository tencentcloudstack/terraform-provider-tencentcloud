package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsWatermarkTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsWatermarkTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_watermark_template.watermark_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_watermark_template.watermark_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsWatermarkTemplate = `

resource "tencentcloud_mps_watermark_template" "watermark_template" {
  type = &lt;nil&gt;
  name = &lt;nil&gt;
  comment = &lt;nil&gt;
  coordinate_origin = "TopLeft"
  x_pos = "0px"
  y_pos = "0px"
  image_template {
		image_content = &lt;nil&gt;
		width = "10%"
		height = "0px"
		repeat_type = "repeat"

  }
  text_template {
		font_type = &lt;nil&gt;
		font_size = &lt;nil&gt;
		font_color = "0xFFFFFF"
		font_alpha = 

  }
  svg_template {
		width = "10W%"
		height = "0px"

  }
}

`
