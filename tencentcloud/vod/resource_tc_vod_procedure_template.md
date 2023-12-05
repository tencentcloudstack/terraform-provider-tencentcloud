Provide a resource to create a VOD procedure template.

Example Usage

```hcl
resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
  drm_type                        = "SimpleAES"
  disable_higher_video_bitrate    = false
  disable_higher_video_resolution = false
  comment                         = "test"

  stream_info {
    video {
      codec               = "libx265"
      fps                 = 4
      bitrate             = 129
      resolution_adaptive = false
      width               = 128
      height              = 128
      fill_type           = "stretch"
    }
    audio {
      codec         = "libmp3lame"
      bitrate       = 129
      sample_rate   = 44100
      audio_channel = "dual"
    }
    remove_audio = false
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
  }
}

resource "tencentcloud_vod_snapshot_by_time_offset_template" "foo" {
  name                = "tf-snapshot"
  width               = 130
  height              = 128
  resolution_adaptive = false
  format              = "png"
  comment             = "test"
  fill_type           = "white"
}

resource "tencentcloud_vod_image_sprite_template" "foo" {
  sample_type         = "Percent"
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

resource "tencentcloud_vod_procedure_template" "foo" {
  name    = "tf-procedure"
  comment = "test"
  media_process_task {
    adaptive_dynamic_streaming_task_list {
      definition = tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id
    }
    snapshot_by_time_offset_task_list {
      definition           = tencentcloud_vod_snapshot_by_time_offset_template.foo.id
      ext_time_offset_list = [
        "3.5s"
      ]
    }
    image_sprite_task_list {
      definition = tencentcloud_vod_image_sprite_template.foo.id
    }
  }
}
```

Import

VOD procedure template can be imported using the name, e.g.

```
$ terraform import tencentcloud_vod_procedure_template.foo tf-procedure
```