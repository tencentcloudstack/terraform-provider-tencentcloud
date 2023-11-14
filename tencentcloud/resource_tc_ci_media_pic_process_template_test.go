package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCiMediaPicProcessTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaPicProcessTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_pic_process_template.media_pic_process_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_pic_process_template.media_pic_process_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaPicProcessTemplate = `

resource "tencentcloud_ci_media_pic_process_template" "media_pic_process_template" {
  name = &lt;nil&gt;
  watermark {
		type = &lt;nil&gt;
		pos = &lt;nil&gt;
		loc_mode = &lt;nil&gt;
		dx = &lt;nil&gt;
		dy = &lt;nil&gt;
		start_time = &lt;nil&gt;
		end_time = &lt;nil&gt;
		slide_config {
			slide_mode = &lt;nil&gt;
			x_slide_speed = &lt;nil&gt;
			x_slide_speed = &lt;nil&gt;
		}
		image {
			url = &lt;nil&gt;
			mode = &lt;nil&gt;
			width = &lt;nil&gt;
			height = &lt;nil&gt;
			transparency = &lt;nil&gt;
			background = &lt;nil&gt;
		}
		text {
			font_size = &lt;nil&gt;
			font_type = &lt;nil&gt;
			font_color = &lt;nil&gt;
			transparency = &lt;nil&gt;
			text = &lt;nil&gt;
		}

  }
}

`
