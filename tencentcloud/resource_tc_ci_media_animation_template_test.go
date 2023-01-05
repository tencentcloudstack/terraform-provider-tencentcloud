package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCiMediaAnimationTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaAnimationTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_animation_template.media_animation_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_animation_template.media_animation_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaAnimationTemplate = `

resource "tencentcloud_ci_media_animation_template" "media_animation_template" {
  name = &lt;nil&gt;
  container {
		format = &lt;nil&gt;

  }
  video {
		codec = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		fps = &lt;nil&gt;
		animate_only_keep_key_frame = &lt;nil&gt;
		animate_time_interval_of_frame = &lt;nil&gt;
		animate_frames_per_second = &lt;nil&gt;
		quality = &lt;nil&gt;

  }
  time_interval {
		start = &lt;nil&gt;
		duration = &lt;nil&gt;

  }
}

`
