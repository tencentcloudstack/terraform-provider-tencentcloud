Provide a resource to create a VOD procedure template.

Example Usage

```hcl
resource  "tencentcloud_vod_sub_application" "sub_application" {
	name = "procedure-subapplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  drm_type                        = "SimpleAES"
  disable_higher_video_bitrate    = false
  disable_higher_video_resolution = false
  comment                         = "test"

  stream_info {
    video {
      codec   = "libx264"
      fps     = 3
      bitrate = 128
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 128
      sample_rate = 32000
    }
    remove_audio = true
  }
  stream_info {
    video {
      codec   = "libx264"
      fps     = 4
      bitrate = 256
    }
    audio {
      codec       = "libfdk_aac"
      bitrate     = 256
      sample_rate = 44100
    }
    remove_audio = true
    tehd_config {
      type = "TEHD-100"
    }
  }
}

resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  width               = 128
  height              = 128
  resolution_adaptive = false
  format              = "png"
  comment             = "test"
  fill_type           = "white"
}

resource "tencentcloud_vod_image_sprite_template" "foo" {
  sample_type         = "Percent"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  sample_interval     = 10
  row_count           = 3
  column_count        = 3
  name                = "tf-sprite"
  comment             = "test"
  fill_type           = "stretch"
  width               = 128
  height              = 128
  resolution_adaptive = false
}

resource "tencentcloud_vod_transcode_template" "transcode_template" {
  container = "mp4"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  name = "720pTranscodeTemplate"
  comment = "test transcode mp4 720p update"
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

resource "tencentcloud_vod_procedure_template" "foo" {
  name    = "tf-procedure0"
  comment = "test"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  media_process_task {
    adaptive_dynamic_streaming_task_list {
      definition = tonumber(split("#", tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id)[1])
    }
    snapshot_by_time_offset_task_list {
      definition = tonumber(split("#", tencentcloud_vod_snapshot_by_time_offset_template.foo.id)[1])
      ext_time_offset_list = [
        "3.5s"
      ]
    }
    image_sprite_task_list {
      definition = tonumber(split("#", tencentcloud_vod_image_sprite_template.foo.id)[1])
    }
    transcode_task_list {
      definition = tonumber(split("#", tencentcloud_vod_transcode_template.transcode_template.id)[1])
    }
  }
}
```

Import

VOD procedure template can be imported using the name, e.g.

```
$ terraform import tencentcloud_vod_procedure_template.foo tf-procedure
```