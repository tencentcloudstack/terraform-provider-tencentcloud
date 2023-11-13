package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCiMediaTranscodeTemplateResource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCiMediaTranscodeTemplate,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttrSet("tencentcloud_ci_media_transcode_template.media_transcode_template", "id")),
			},
			{
				ResourceName:      "tencentcloud_ci_media_transcode_template.media_transcode_template",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

const testAccCiMediaTranscodeTemplate = `

resource "tencentcloud_ci_media_transcode_template" "media_transcode_template" {
  name = &lt;nil&gt;
  container {
		format = &lt;nil&gt;
		clip_config {
			duration = &lt;nil&gt;
		}

  }
  video {
		codec = &lt;nil&gt;
		width = &lt;nil&gt;
		height = &lt;nil&gt;
		fps = &lt;nil&gt;
		remove = &lt;nil&gt;
		profile = &lt;nil&gt;
		bitrate = &lt;nil&gt;
		crf = &lt;nil&gt;
		gop = &lt;nil&gt;
		preset = &lt;nil&gt;
		bufsize = &lt;nil&gt;
		maxrate = &lt;nil&gt;
		pixfmt = &lt;nil&gt;
		long_short_mode = &lt;nil&gt;
		rotate = &lt;nil&gt;

  }
  time_interval {
		start = &lt;nil&gt;
		duration = &lt;nil&gt;

  }
  audio {
		codec = &lt;nil&gt;
		samplerate = &lt;nil&gt;
		bitrate = &lt;nil&gt;
		channels = &lt;nil&gt;
		remove = &lt;nil&gt;
		keep_two_tracks = &lt;nil&gt;
		switch_track = &lt;nil&gt;
		sample_format = &lt;nil&gt;

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
		transcode_index = &lt;nil&gt;
		hls_encrypt {
			is_hls_encrypt = &lt;nil&gt;
			uri_key = &lt;nil&gt;
		}
		dash_encrypt {
			is_encrypt = &lt;nil&gt;
			uri_key = &lt;nil&gt;
		}

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
}

`
