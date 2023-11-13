package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsTranscodeTemplatesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsTranscodeTemplatesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mps_transcode_templates.transcode_templates")),
			},
		},
	})
}

const testAccMpsTranscodeTemplatesDataSource = `

data "tencentcloud_mps_transcode_templates" "transcode_templates" {
  definitions = &lt;nil&gt;
  type = &lt;nil&gt;
  container_type = &lt;nil&gt;
  t_e_h_d_type = &lt;nil&gt;
  offset = &lt;nil&gt;
  limit = &lt;nil&gt;
  transcode_type = &lt;nil&gt;
  total_count = &lt;nil&gt;
  transcode_template_set {
		definition = &lt;nil&gt;
		container = &lt;nil&gt;
		name = &lt;nil&gt;
		comment = &lt;nil&gt;
		type = &lt;nil&gt;
		remove_video = &lt;nil&gt;
		remove_audio = &lt;nil&gt;
		video_template {
			codec = &lt;nil&gt;
			fps = &lt;nil&gt;
			bitrate = &lt;nil&gt;
			resolution_adaptive = &lt;nil&gt;
			width = &lt;nil&gt;
			height = &lt;nil&gt;
			gop = &lt;nil&gt;
			fill_type = &lt;nil&gt;
			vcrf = &lt;nil&gt;
		}
		audio_template {
			codec = &lt;nil&gt;
			bitrate = &lt;nil&gt;
			sample_rate = &lt;nil&gt;
			audio_channel = &lt;nil&gt;
		}
		t_e_h_d_config {
			type = &lt;nil&gt;
			max_video_bitrate = &lt;nil&gt;
		}
		container_type = &lt;nil&gt;
		create_time = &lt;nil&gt;
		update_time = &lt;nil&gt;
		enhance_config {
			video_enhance {
				frame_rate {
					switch = &lt;nil&gt;
					fps = &lt;nil&gt;
				}
				super_resolution {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
					size = &lt;nil&gt;
				}
				hdr {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
				}
				denoise {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
				}
				image_quality_enhance {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
				}
				color_enhance {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
				}
				sharp_enhance {
					switch = &lt;nil&gt;
					intensity = &lt;nil&gt;
				}
				face_enhance {
					switch = &lt;nil&gt;
					intensity = &lt;nil&gt;
				}
				low_light_enhance {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
				}
				scratch_repair {
					switch = &lt;nil&gt;
					intensity = &lt;nil&gt;
				}
				artifact_repair {
					switch = &lt;nil&gt;
					type = &lt;nil&gt;
				}
			}
		}

  }
}

`
