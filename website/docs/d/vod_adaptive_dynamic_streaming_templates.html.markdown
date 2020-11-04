---
subcategory: "VOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_adaptive_dynamic_streaming_templates"
sidebar_current: "docs-tencentcloud-datasource-vod_adaptive_dynamic_streaming_templates"
description: |-
  Use this data source to query detailed information of VOD adaptive dynamic streaming templates.
---

# tencentcloud_vod_adaptive_dynamic_streaming_templates

Use this data source to query detailed information of VOD adaptive dynamic streaming templates.

## Example Usage

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

data "tencentcloud_vod_adaptive_dynamic_streaming_templates" "foo" {
  type       = "Custom"
  definition = tencentcloud_vod_adaptive_dynamic_streaming_template.foo.id
}
```

## Argument Reference

The following arguments are supported:

* `definition` - (Optional) Unique ID filter of adaptive dynamic streaming template.
* `result_output_file` - (Optional) Used to save results.
* `sub_app_id` - (Optional) Subapplication ID in VOD. If you need to access a resource in a subapplication, enter the subapplication ID in this field; otherwise, leave it empty.
* `type` - (Optional) Template type filter. Valid values: `Preset`, `Custom`. `Preset`: preset template; `Custom`: custom template.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `template_list` - A list of adaptive dynamic streaming templates. Each element contains the following attributes:
  * `comment` - Template description.
  * `create_time` - Creation time of template in ISO date format.
  * `definition` - Unique ID of adaptive dynamic streaming template.
  * `disable_higher_video_bitrate` - Whether to prohibit transcoding video from low bitrate to high bitrate. `false`: no, `true`: yes.
  * `disable_higher_video_resolution` - Whether to prohibit transcoding from low resolution to high resolution. `false`: no, `true`: yes.
  * `drm_type` - DRM scheme type.
  * `format` - Adaptive bitstream format.
  * `name` - Template name.
  * `stream_info` - List of AdaptiveStreamTemplate parameter information of output substream for adaptive bitrate streaming.
    * `audio` - Audio parameter information.
      * `audio_channel` - Audio channel system. Valid values: mono, dual, stereo.
      * `bitrate` - Audio stream bitrate in Kbps. Value range: `0` and `[26, 256]`. If the value is `0`, the bitrate of the audio stream will be the same as that of the original audio.
      * `codec` - Audio stream encoder. Valid value are: `libfdk_aac` and `libmp3lame`.
      * `sample_rate` - Audio stream sample rate. Valid values: `32000`, `44100`, `48000`. Unit is HZ.
    * `remove_audio` - Whether to remove audio stream. `false`: no, `true`: yes.
    * `video` - Video parameter information.
      * `bitrate` - Bitrate of video stream in Kbps. Value range: `0` and `[128, 35000]`. If the value is `0`, the bitrate of the video will be the same as that of the source video.
      * `codec` - Video stream encoder. Valid values: `libx264`, `libx265`, `av1`.`libx264`: H.264, `libx265`: H.265, `av1`: AOMedia Video 1. Currently, a resolution within 640x480 must be specified for `H.265`. and the `av1` container only supports mp4.
      * `fill_type` - Fill type. Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported: `stretch`: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot shorter or longer; `black`: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks. Note: this field may return null, indicating that no valid values can be obtained.
      * `fps` - Video frame rate in Hz. Value range: `[0, 60]`. If the value is `0`, the frame rate will be the same as that of the source video.
      * `height` - Maximum value of the height (or short side) of a video stream in px. Value range: `0` and `[128, 4096]`. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Note: this field may return null, indicating that no valid values can be obtained.
      * `resolution_adaptive` - Resolution adaption. Valid values: `true`,`false`. `true`: enabled. In this case, `width` represents the long side of a video, while `height` the short side; `false`: disabled. In this case, `width` represents the width of a video, while `height` the height. Note: this field may return null, indicating that no valid values can be obtained.
      * `width` - Maximum value of the width (or long side) of a video stream in px. Value range: `0` and `[128, 4096]`. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Note: this field may return null, indicating that no valid values can be obtained.
  * `type` - Template type filter. Valid values: `Preset`,`Custom`. `Preset`: preset template; `Custom`: custom template.
  * `update_time` - Last modified time of template in ISO date format.


