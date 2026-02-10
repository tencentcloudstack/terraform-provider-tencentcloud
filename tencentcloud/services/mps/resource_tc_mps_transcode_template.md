Provides a resource to create a MPS transcode template

Example Usage

```hcl
resource "tencentcloud_mps_transcode_template" "example" {
  container    = "mp4"
  name         = "tf-example"
  remove_audio = 0
  remove_video = 0

  audio_template {
    audio_channel = 2
    bitrate       = 27
    codec         = "libfdk_aac"
    sample_rate   = 32000
  }

  video_template {
    bitrate             = 130
    codec               = "libx264"
    fill_type           = "black"
    fps                 = 20
    gop                 = 0
    height              = 4096
    resolution_adaptive = "close"
    vcrf                = 0
    width               = 128
  }
}
```

Import

MPS transcode template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_transcode_template.example 1638607
```