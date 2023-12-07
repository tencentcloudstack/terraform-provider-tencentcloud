Provides a resource to create a css live_transcode_template

Example Usage

```hcl
resource "tencentcloud_css_live_transcode_template" "live_transcode_template" {
  template_name = "template_name"
  acodec = "aac"
  audio_bitrate = 128
  video_bitrate = 100
  vcodec = "origin"
  description = "This_is_a_tf_test_temp."
  need_video = 1
  width = 0
  need_audio = 1
  height = 0
  fps = 0
  gop = 2
  rotate = 0
  profile = "baseline"
  bitrate_to_orig = 0
  height_to_orig = 0
  fps_to_orig = 0
  ai_trans_code = 0
  adapt_bitrate_percent = 0
  short_edge_as_height = 0
  drm_type = "fairplay"
  drm_tracks = "SD"
}

```
Import

css live_transcode_template can be imported using the id, e.g.
```
$ terraform import tencentcloud_css_live_transcode_template.live_transcode_template liveTranscodeTemplate_id
```