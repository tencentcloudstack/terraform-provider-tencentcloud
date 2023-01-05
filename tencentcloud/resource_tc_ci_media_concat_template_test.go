package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCiMediaConcatTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaConcatTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_concat_template.media_concat_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_concat_template.media_concat_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaConcatTemplate = `

resource "tencentcloud_ci_media_concat_template" "media_concat_template" {
  name = &lt;nil&gt;
  concat_template {
		concat_fragment {
			url = &lt;nil&gt;
			mode = &lt;nil&gt;
		}
		audio {
			codec = &lt;nil&gt;
			samplerate = &lt;nil&gt;
			bitrate = &lt;nil&gt;
			channels = &lt;nil&gt;
		}
		video {
			codec = &lt;nil&gt;
			width = &lt;nil&gt;
			height = &lt;nil&gt;
			fps = &lt;nil&gt;
			crf = &lt;nil&gt;
			remove = &lt;nil&gt;
			rotate = &lt;nil&gt;
		}
		container {
			format = &lt;nil&gt;
		}
		audio_mix {
			audio_source = &lt;nil&gt;
			mix_mode = &lt;nil&gt;
			replace = &lt;nil&gt;
			effect_config {
				enable_start_fadein = &lt;nil&gt;
				start_fadein_time = &lt;nil&gt;
				enable_end_fadeout = &lt;nil&gt;
				end_fadeout_time = &lt;nil&gt;
				enable_bgm_fade = &lt;nil&gt;
				bgm_fade_time = &lt;nil&gt;
			}
		}
		direct_concat = &lt;nil&gt;

  }
}

`
