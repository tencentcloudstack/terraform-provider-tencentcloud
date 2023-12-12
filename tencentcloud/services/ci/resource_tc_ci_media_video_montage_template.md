Provides a resource to create a ci media_video_montage_template

Example Usage

```hcl
resource "tencentcloud_ci_media_video_montage_template" "media_video_montage_template" {
  bucket = "terraform-ci-xxxxx"
  name = "video_montage_template"
  duration = "10.5"
  audio {
		codec = "aac"
		samplerate = "44100"
		bitrate = "128"
		channels = "4"
		remove = "false"

  }
  video {
		codec = "H.264"
		width = "1280"
		height = ""
		bitrate = "1000"
		fps = "25"
		crf = ""
		remove = ""
  }
  container {
		format = "mp4"

  }
  audio_mix {
		audio_source = "https://terraform-ci-xxxxx.cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"
		mix_mode = "Once"
		replace = "true"
		# effect_config {
		# 	enable_start_fadein = ""
		# 	start_fadein_time = ""
		# 	enable_end_fadeout = ""
		# 	end_fadeout_time = ""
		# 	enable_bgm_fade = ""
		# 	bgm_fade_time = ""
		# }

  }
}
```

Import

ci media_video_montage_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_video_montage_template.media_video_montage_template terraform-ci-xxxxxx#t193e5ecc1b8154e57a8376b4405ad9c63
```