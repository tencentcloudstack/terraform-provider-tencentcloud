Provides a resource to create a TEO just-in-time transcode template.

Example Usage

```hcl
resource "tencentcloud_teo_just_in_time_transcode_template" "example" {
  zone_id       = "zone-2qtuhspy7cr6"
  template_name = "my-template"
  comment       = "Example transcode template"

  video_stream_switch = "on"
  audio_stream_switch = "on"

  video_template {
    video_codec         = "H.264"
    fps                 = 30
    bitrate             = 2000
    resolution_adaptive = "open"
    width               = 1280
    height              = 720
    fill_type           = "black"
  }

  audio_template {
    codec         = "libfdk_aac"
    audio_channel = 2
  }
}
```

Import

TEO just-in-time transcode template can be imported using the joint id "zone_id#template_id", e.g.

```
terraform import tencentcloud_teo_just_in_time_transcode_template.example zone-2qtuhspy7cr6#jitt-abcdefghij
```
