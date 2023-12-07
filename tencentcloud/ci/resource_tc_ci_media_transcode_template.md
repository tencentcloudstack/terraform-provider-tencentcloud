Provides a resource to create a ci media_transcode_template

Example Usage

```hcl
resource "tencentcloud_ci_media_transcode_template" "media_transcode_template" {
  bucket = "terraform-ci-1308919341"
  name = "transcode_template"
  container {
		format = "mp4"
		# clip_config {
		# 	duration = ""
		# }
  }
  video {
		codec = "H.264"
		width = "1280"
		# height = ""
		fps = "30"
		remove = "false"
		profile = "high"
		bitrate = "1000"
		# crf = ""
		# gop = ""
		preset = "medium"
		# bufsize = ""
		# maxrate = ""
		# pixfmt = ""
		long_short_mode = "false"
		# rotate = ""
  }
  time_interval {
		start = "0"
		duration = "60"
  }
  audio {
		codec = "aac"
		samplerate = "44100"
		bitrate = "128"
		channels = "4"
		remove = "false"
		keep_two_tracks = "false"
		switch_track = "false"
		sample_format = ""
  }
  trans_config {
		adj_dar_method = "scale"
		is_check_reso = "false"
		reso_adj_method = "1"
		is_check_video_bitrate = "false"
		video_bitrate_adj_method = "0"
		is_check_audio_bitrate = "false"
		audio_bitrate_adj_method = "0"
		delete_metadata = "false"
		is_hdr2_sdr = "false"
  }
  audio_mix {
		audio_source = "https://terraform-ci-1308919341.cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"
		mix_mode = "Once"
		replace = "true"
		effect_config {
			enable_start_fadein = "true"
			start_fadein_time = "3"
			enable_end_fadeout = "false"
			end_fadeout_time = "0"
			enable_bgm_fade = "true"
			bgm_fade_time = "1.7"
		}
  }
}
```

Import

ci media_transcode_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_transcode_template.media_transcode_template media_transcode_template_id
```