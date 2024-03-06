Provides a resource to create a vod transcode template

Example Usage

```hcl
resource  "tencentcloud_vod_sub_application" "sub_application" {
	name = "transcodeTemplateSubApplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_transcode_template" "transcode_template" {
  container = "mp4"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  name = "720pTranscodeTemplate"
  comment = "test transcode mp4 720p"
  remove_video = 0
  remove_audio = 0
  video_template {
	codec = "libx264"
	fps = 26
	bitrate = 1000
	resolution_adaptive = "open"
	width = 0
	height = 720
	fill_type = "stretch"
	vcrf = 1
	gop = 250
	preserve_hdr_switch = "OFF"
	codec_tag = "hvc1"

  }
  audio_template {
	codec = "libfdk_aac"
	bitrate = 128
	sample_rate = 44100
	audio_channel = 2
  }
  segment_type = "ts"
}
```

Import

vod transcode template can be imported using the id, e.g.

```
terraform import tencentcloud_vod_transcode_template.transcode_template $subAppId#$templateId
```