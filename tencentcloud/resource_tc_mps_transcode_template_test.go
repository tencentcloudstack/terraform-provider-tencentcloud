package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMpsTranscodeTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMpsTranscodeTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_mps_transcode_template.transcode_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_mps_transcode_template.transcode_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccMpsTranscodeTemplate = `

resource "tencentcloud_mps_transcode_template" "transcode_template" {
  container = &lt;nil&gt;
  name = &lt;nil&gt;
  comment = &lt;nil&gt;
  remove_video = 0
  remove_audio = 0
  video_template {
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
  audio_template {
		codec = &lt;nil&gt;
		bitrate = &lt;nil&gt;
		sample_rate = &lt;nil&gt;
		audio_channel = 2

  }
  t_e_h_d_config {
		type = &lt;nil&gt;
		max_video_bitrate = &lt;nil&gt;

  }
  enhance_config {
		video_enhance {
			frame_rate {
				switch = "ON"
				fps = 0
			}
			super_resolution {
				switch = "ON"
				type = "lq"
				size = 2
			}
			hdr {
				switch = "ON"
				type = "HDR10"
			}
			denoise {
				switch = "ON"
				type = "weak"
			}
			image_quality_enhance {
				switch = "ON"
				type = "weak"
			}
			color_enhance {
				switch = "ON"
				type = "weak"
			}
			sharp_enhance {
				switch = "ON"
				intensity = 
			}
			face_enhance {
				switch = "ON"
				intensity = 
			}
			low_light_enhance {
				switch = "ON"
				type = "normal"
			}
			scratch_repair {
				switch = "ON"
				intensity = 
			}
			artifact_repair {
				switch = "ON"
				type = "weak"
			}
		}

  }
}

`
