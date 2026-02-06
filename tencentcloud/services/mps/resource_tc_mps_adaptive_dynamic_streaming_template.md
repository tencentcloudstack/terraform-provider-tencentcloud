Provides a resource to create a MPS adaptive dynamic streaming template

Example Usage

```hcl
resource "tencentcloud_mps_adaptive_dynamic_streaming_template" "example" {
  name                            = "tf-example"
  comment                         = "terrraform test"
  disable_higher_video_bitrate    = 0
  disable_higher_video_resolution = 1
  format                          = "HLS"
  pure_audio                      = 0
  segment_type                    = "ts-segment"
  stream_infos {
    remove_audio = 0
    remove_video = 0

    audio {
      audio_channel = 1
      bitrate       = 55
      codec         = "libmp3lame"
      sample_rate   = 32000
    }

    video {
      bitrate             = 245
      codec               = "libx264"
      fill_type           = "black"
      fps                 = 30
      gop                 = 0
      height              = 135
      resolution_adaptive = "open"
      vcrf                = 0
      width               = 145
    }
  }

  stream_infos {
    remove_audio = 0
    remove_video = 0

    audio {
      audio_channel = 2
      bitrate       = 60
      codec         = "libfdk_aac"
      sample_rate   = 32000
    }

    video {
      bitrate             = 400
      codec               = "libx264"
      fill_type           = "black"
      fps                 = 40
      gop                 = 0
      height              = 150
      resolution_adaptive = "open"
      vcrf                = 0
      width               = 160
    }
  }
}
```

Import

MPS adaptive dynamic streaming template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_adaptive_dynamic_streaming_template.example 1636009
```