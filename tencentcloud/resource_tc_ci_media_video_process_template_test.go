package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCiMediaVideoProcessTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaVideoProcessTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_video_process_template.media_video_process_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_video_process_template.media_video_process_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaVideoProcessTemplate = `

resource "tencentcloud_ci_media_video_process_template" "media_video_process_template" {
  name = &lt;nil&gt;
  color_enhance {
		enable = &lt;nil&gt;
		contrast = &lt;nil&gt;
		correction = &lt;nil&gt;
		saturation = &lt;nil&gt;

  }
  ms_sharpen {
		enable = &lt;nil&gt;
		sharpen_level = &lt;nil&gt;

  }
}

`
