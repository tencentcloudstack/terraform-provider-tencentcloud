package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
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
  name = ""
  container {
		format = ""
		clip_config {
			duration = ""
		}

  }
  video {
		codec = ""
		width = ""
		height = ""
		fps = ""
		remove = ""
		profile = ""
		bitrate = ""
		crf = ""
		gop = ""
		preset = ""
		bufsize = ""
		maxrate = ""
		pixfmt = ""
		long_short_mode = ""
		rotate = ""

  }
  time_interval {
		start = ""
		duration = ""

  }
  audio {
		codec = ""
		samplerate = ""
		bitrate = ""
		bitrate = ""
		remove = ""
		keep_two_tracks = ""
		switch_track = ""
		sample_format = ""

  }
  trans_config {
		adj_dar_method = ""
		is_check_reso = ""
		reso_adj_method = ""
		is_check_video_bitrate = ""
		video_bitrate_adj_method = ""
		is_check_audio_bitrate = ""
		audio_bitrate_adj_method = ""
		is_check_video_fps = ""
		video_fps_adj_method = ""
		delete_metadata = ""
		is_hdr2_sdr = ""
		transcode_index = ""
		hls_encrypt {
			is_hls_encrypt = ""
			uri_key = ""
		}
		dash_encrypt {
			is_encrypt = ""
			uri_key = ""
		}

  }
  audio_mix {
		audio_source = ""
		mix_mode = ""
		replace = ""
		effect_config {
			enable_start_fadein = ""
			start_fadein_time = ""
			enable_end_fadeout = ""
			end_fadeout_time = ""
			enable_bgm_fade = ""
			bgm_fade_time = ""
		}

  }
  audio_mix_array {
		audio_source = ""
		mix_mode = ""
		replace = ""
		effect_config {
			enable_start_fadein = ""
			start_fadein_time = ""
			enable_end_fadeout = ""
			end_fadeout_time = ""
			enable_bgm_fade = ""
			bgm_fade_time = ""
		}

  }
}

`
