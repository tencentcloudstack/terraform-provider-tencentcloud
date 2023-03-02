package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudMpsAdaptiveDynamicStreamingTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsAdaptiveDynamicStreamingTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_adaptive_dynamic_streaming_template.adaptive_dynamic_streaming_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_adaptive_dynamic_streaming_template.adaptive_dynamic_streaming_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsAdaptiveDynamicStreamingTemplate = `

resource "tencentcloud_mps_adaptive_dynamic_streaming_template" "adaptive_dynamic_streaming_template" {
  format = &lt;nil&gt;
  stream_infos {
		video {
			codec = &lt;nil&gt;
			fps = &lt;nil&gt;
			bitrate = &lt;nil&gt;
			resolution_adaptive = "open"
			width = 0
			height = 0
			gop = &lt;nil&gt;
			fill_type = "black"
			vcrf = &lt;nil&gt;
		}
		audio {
			codec = &lt;nil&gt;
			bitrate = &lt;nil&gt;
			sample_rate = &lt;nil&gt;
			audio_channel = 2
		}
		remove_audio = &lt;nil&gt;
		remove_video = &lt;nil&gt;

  }
  name = &lt;nil&gt;
  disable_higher_video_bitrate = 0
  disable_higher_video_resolution = 0
  comment = &lt;nil&gt;
}

`
