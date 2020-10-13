---
subcategory: "VOD"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_adaptive_dynamic_streaming_template"
sidebar_current: "docs-tencentcloud-resource-vod_adaptive_dynamic_streaming_template"
description: |-
  Provide a resource to create a VOD adaptive dynamic streaming template.
---

# tencentcloud_vod_adaptive_dynamic_streaming_template

Provide a resource to create a VOD adaptive dynamic streaming template.

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
  }
}
```

## Argument Reference

The following arguments are supported:

* `format` - (Required) Adaptive bitstream format. Valid values: `HLS`.
* `name` - (Required) Template name. Length limit: 64 characters.
* `stream_info` - (Required) List of AdaptiveStreamTemplate parameter information of output substream for adaptive bitrate streaming. Up to 10 substreams can be output. Note: the frame rate of all substreams must be the same; otherwise, the frame rate of the first substream will be used as the output frame rate.
* `comment` - (Optional) Template description. Length limit: 256 characters.
* `disable_higher_video_bitrate` - (Optional) Whether to prohibit transcoding video from low bitrate to high bitrate. Valid values: `false`: no, `true`: yes. Default value: `false`.
* `disable_higher_video_resolution` - (Optional) Whether to prohibit transcoding from low resolution to high resolution. Valid values: `false`: no, `true`: yes. Default value: `false`.
* `drm_type` - (Optional, ForceNew) DRM scheme type. Valid values: `SimpleAES`. If this field is an empty string, DRM will not be performed on the video.
* `sub_app_id` - (Optional) Subapplication ID in VOD. If you need to access a resource in a subapplication, enter the subapplication ID in this field; otherwise, leave it empty.

The `audio` object supports the following:

* `bitrate` - (Required) Audio stream bitrate in Kbps. Value range: `0` and `[26, 256]`. If the value is `0`, the bitrate of the audio stream will be the same as that of the original audio.
* `codec` - (Required) Audio stream encoder. Valid value are: `libfdk_aac` and `libmp3lame`, while `libfdk_aac` is recommended.
* `sample_rate` - (Required) Audio stream sample rate. Valid values: `32000`, `44100`, `48000`, in Hz.
* `audio_channel` - (Optional) Audio channel system. Valid values: mono, dual, stereo. Default value: dual.

The `stream_info` object supports the following:

* `audio` - (Required) Audio parameter information.
* `video` - (Required) Video parameter information.
* `remove_audio` - (Optional) Whether to remove audio stream. Valid values: `false`: no, `true`: yes. `false` by default.

The `video` object supports the following:

* `bitrate` - (Required) Bitrate of video stream in Kbps. Value range: `0` and `[128, 35000]`. If the value is `0`, the bitrate of the video will be the same as that of the source video.
* `codec` - (Required) Video stream encoder. Valid values: `libx264`: H.264, `libx265`: H.265, `av1`: AOMedia Video 1. Currently, a resolution within 640x480 must be specified for `H.265`. and the `av1` container only supports mp4.
* `fps` - (Required) Video frame rate in Hz. Value range: `[0, 60]`. If the value is `0`, the frame rate will be the same as that of the source video.
* `fill_type` - (Optional) Fill type. Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported: `stretch`: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot shorter or longer; `black`: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks. Default value: black. Note: this field may return null, indicating that no valid values can be obtained.
* `height` - (Optional) Maximum value of the height (or short side) of a video stream in px. Value range: `0` and `[128, 4096]`. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Default value: `0`. Note: this field may return null, indicating that no valid values can be obtained.
* `resolution_adaptive` - (Optional) Resolution adaption. Valid values: `true`: enabled. In this case, `width` represents the long side of a video, while `height` the short side; `false`: disabled. In this case, `width` represents the width of a video, while `height` the height. Default value: `true`. Note: this field may return null, indicating that no valid values can be obtained.
* `width` - (Optional) Maximum value of the width (or long side) of a video stream in px. Value range: `0` and `[128, 4096]`. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Default value: `0`. Note: this field may return null, indicating that no valid values can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of template in ISO date format.
* `update_time` - Last modified time of template in ISO date format.


## Import

Vod adaptive dynamic streaming template can be imported using the id, e.g.

```
$ terraform import tencentcloud_vod_adaptive_dynamic_streaming_template.foo 169141
```

