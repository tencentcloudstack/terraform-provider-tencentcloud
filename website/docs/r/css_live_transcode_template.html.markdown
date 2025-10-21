---
subcategory: "Cloud Streaming Services(CSS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_css_live_transcode_template"
sidebar_current: "docs-tencentcloud-resource-css_live_transcode_template"
description: |-
  Provides a resource to create a css live_transcode_template
---

# tencentcloud_css_live_transcode_template

Provides a resource to create a css live_transcode_template

## Example Usage

```hcl
resource "tencentcloud_css_live_transcode_template" "live_transcode_template" {
  template_name         = "template_name"
  acodec                = "aac"
  audio_bitrate         = 128
  video_bitrate         = 100
  vcodec                = "origin"
  description           = "This_is_a_tf_test_temp."
  need_video            = 1
  width                 = 0
  need_audio            = 1
  height                = 0
  fps                   = 0
  gop                   = 2
  rotate                = 0
  profile               = "baseline"
  bitrate_to_orig       = 0
  height_to_orig        = 0
  fps_to_orig           = 0
  ai_trans_code         = 0
  adapt_bitrate_percent = 0
  short_edge_as_height  = 0
  drm_type              = "fairplay"
  drm_tracks            = "SD"
}
```

## Argument Reference

The following arguments are supported:

* `template_name` - (Required, String) template name, only support 0-9 and a-z.
* `video_bitrate` - (Required, Int) video bitrate, 0 for origin, range 0kbps - 8000kbps.
* `acodec` - (Optional, String) default aac, not support now.
* `adapt_bitrate_percent` - (Optional, Float64) high speed mode adapt bitrate, support 0 - 0.5.
* `ai_trans_code` - (Optional, Int) enable high speed mode, default 0, 1 for enable, 0 for no.
* `audio_bitrate` - (Optional, Int) default 0, range 0 - 500.
* `bitrate_to_orig` - (Optional, Int) base on origin bitrate if origin bitrate is lower than the setting bitrate. default 0, 1 for yes, 0 for no.
* `description` - (Optional, String) template desc.
* `drm_tracks` - (Optional, String) DRM tracks, support AUDIO/SD/HD/UHD1/UHD2.
* `drm_type` - (Optional, String) DRM type, support fairplay/normalaes/widevine.
* `fps_to_orig` - (Optional, Int) base on origin fps if origin fps is lower than the setting fps. default 0, 1 for yes, 0 for no.
* `fps` - (Optional, Int) video fps, default 0, range 0 - 60.
* `gop` - (Optional, Int) gop of the video, second, default origin of the video, range 2 - 6.
* `height_to_orig` - (Optional, Int) base on origin height if origin height is lower than the setting height. default 0, 1 for yes, 0 for no.
* `height` - (Optional, Int) template height, default 0, range 0 - 3000, must be pow of 2, needed while AiTransCode = 1.
* `need_audio` - (Optional, Int) keep audio or not, default 1 for yes, 0 for no.
* `need_video` - (Optional, Int) keep video or not, default 1 for yes, 0 for no.
* `profile` - (Optional, String) quality of the video, default baseline, support baseline/main/high.
* `rotate` - (Optional, Int) roate degree, default 0, support 0/90/180/270.
* `short_edge_as_height` - (Optional, Int) let the short edge as the height.
* `vcodec` - (Optional, String) video codec, default origin, support h264/h265/origin.
* `width` - (Optional, Int) template width, default 0, range 0 - 3000, must be pow of 2.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

css live_transcode_template can be imported using the id, e.g.
```
$ terraform import tencentcloud_css_live_transcode_template.live_transcode_template liveTranscodeTemplate_id
```

