package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsWatermarkTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsWatermarkTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_watermark_templates.watermark_templates")),
			},
		},
	})
}

const testAccMpsWatermarkTemplatesDataSource = `

data "tencentcloud_mps_watermark_templates" "watermark_templates" {
  definitions = &lt;nil&gt;
  type = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  total_count = &lt;nil&gt;
  watermark_template_set {
		definition = &lt;nil&gt;
		type = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		x_pos = &lt;nil&gt;
		y_pos = &lt;nil&gt;
		image_template {
			image_url = &lt;nil&gt;
			width = &lt;nil&gt;
			height = &lt;nil&gt;
			repeat_type = &lt;nil&gt;
		}
		text_template {
			font_type = &lt;nil&gt;
			font_size = &lt;nil&gt;
			font_color = &lt;nil&gt;
			font_alpha = &lt;nil&gt;
		}
		svg_template {
			width = &lt;nil&gt;
			height = &lt;nil&gt;
		}
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;
		coordinate_origin = &lt;nil&gt;

  }
}

`
