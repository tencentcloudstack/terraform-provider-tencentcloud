package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsAdaptiveDynamicStreamingTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsAdaptiveDynamicStreamingTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_adaptive_dynamic_streaming_templates.adaptive_dynamic_streaming_templates")),
			},
		},
	})
}

const testAccMpsAdaptiveDynamicStreamingTemplatesDataSource = `

data "tencentcloud_mps_adaptive_dynamic_streaming_templates" "adaptive_dynamic_streaming_templates" {
  definitions = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  adaptive_dynamic_streaming_template_set {
		definition = &lt;nil&gt;
		type = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
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
		disable_higher_video_bitrate = &lt;nil&gt;
		disable_higher_video_resolution = &lt;nil&gt;
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;

  }
}

`
