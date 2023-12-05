Provides a resource to create a ci media_concat_template

Example Usage

```hcl
resource "tencentcloud_ci_media_concat_template" "media_concat_template" {
  bucket = "terraform-ci-xxxxxx"
  name = "concat_templates"
  concat_template {
		concat_fragment {
			url = "https://terraform-ci-xxxxxx.cos.ap-guangzhou.myqcloud.com/mp4%2Fmp4-test.mp4"
			mode = "Start"
		}
    concat_fragment {
			url = "https://terraform-ci-xxxxxx.cos.ap-guangzhou.myqcloud.com/mp4%2Fmp4-test.mp4"
			mode = "End"
		}
		audio {
			codec = "mp3"
			samplerate = ""
			bitrate = ""
			channels = ""
		}
		video {
			codec = "H.264"
			width = "1280"
			height = ""
      		bitrate = "1000"
			fps = "25"
			crf = ""
			remove = ""
			rotate = ""
		}
		container {
			format = "mp4"
		}
		audio_mix {
			audio_source = "https://terraform-ci-xxxxxx.cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"
			mix_mode = "Once"
			replace = "true"
			effect_config {
				enable_start_fadein = "true"
				start_fadein_time = "3"
				enable_end_fadeout = "false"
				end_fadeout_time = "0.1"
				enable_bgm_fade = "true"
				bgm_fade_time = "1.7"
			}
		}
  }
}
```

Import

ci media_concat_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_concat_template.media_concat_template id=terraform-ci-xxxxxx#t1cb115dfa1fcc414284f83b7c69bcedcf
```