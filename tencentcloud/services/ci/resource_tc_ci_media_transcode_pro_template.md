Provides a resource to create a ci media_transcode_pro_template

Example Usage

```hcl
resource "tencentcloud_ci_media_transcode_pro_template" "media_transcode_pro_template" {
  bucket = "terraform-ci-xxxxxx"
  name = "transcode_pro_template"
  container {
		format = "mxf"
		# clip_config {
		# 	duration = ""
		# }

  }
  video {
		codec = "xavc"
		profile = "XAVC-HD_422_10bit"
		width = "1920"
		height = "1080"
    	interlaced = "true"
		fps = "30000/1001"
		bitrate = "50000"
		# rotate = ""

  }
  time_interval {
		start = ""
		duration = ""

  }
  audio {
		codec = "pcm_s24le"
		remove = "true"

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
}
```

Import

ci media_transcode_pro_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template terraform-ci-xxxxxx#t13ed9af009da0414e9c7c63456ec8f4d2
```