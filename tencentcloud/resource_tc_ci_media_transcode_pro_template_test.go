package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTencentCloudCiMediaTranscodeProTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaTranscodeProTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaTranscodeProTemplate = `

resource "tencentcloud_ci_media_transcode_pro_template" "media_transcode_pro_template" {
  name = &lt;nil&gt;
  container {
		format = &lt;nil&gt;
		clip_config {
			duration = &lt;nil&gt;
		}

  }
  video {
		codec = &lt;nil&gt;
		profile = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		interlaced = &lt;nil&gt;
		fps = &lt;nil&gt;
		bitrate = &lt;nil&gt;
		rotate = &lt;nil&gt;

  }
  time_interval {
		start = &lt;nil&gt;
		duration = &lt;nil&gt;

  }
  audio {
		codec = &lt;nil&gt;
		remove = &lt;nil&gt;

  }
  trans_config {
		adj_dar_method = &lt;nil&gt;
		is_check_reso = &lt;nil&gt;
		reso_adj_method = &lt;nil&gt;
		is_check_video_bitrate = &lt;nil&gt;
		video_bitrate_adj_method = &lt;nil&gt;
		is_check_audio_bitrate = &lt;nil&gt;
		audio_bitrate_adj_method = &lt;nil&gt;
		is_check_video_fps = &lt;nil&gt;
		video_fps_adj_method = &lt;nil&gt;
		delete_metadata = &lt;nil&gt;
		is_hdr2_sdr = &lt;nil&gt;

  }
}

`
