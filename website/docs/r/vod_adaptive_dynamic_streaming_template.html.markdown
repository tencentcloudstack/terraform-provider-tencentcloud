---
subcategory: "Video on Demand(VOD)"
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
resource "tencentcloud_vod_sub_application" "sub_application" {
  name        = "adaptive-subapplication"
  status      = "On"
  description = "this is sub application"
}

resource "tencentcloud_vod_adaptive_dynamic_streaming_template" "foo" {
  format                          = "HLS"
  name                            = "tf-adaptive"
  sub_app_id                      = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
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
    tehd_config {
      type = "TEHD-100"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `format` - (Required, String) Adaptive bitstream format. Valid values: `HLS`.
* `name` - (Required, String) Template name. Length limit: 64 characters.
* `stream_info` - (Required, List) List of AdaptiveStreamTemplate parameter information of output substream for adaptive bitrate streaming. Up to 10 substreams can be output. Note: the frame rate of all substreams must be the same; otherwise, the frame rate of the first substream will be used as the output frame rate.
* `comment` - (Optional, String) Template description. Length limit: 256 characters.
* `disable_higher_video_bitrate` - (Optional, Bool) Whether to prohibit transcoding video from low bitrate to high bitrate. Valid values: `false`,`true`. `false`: no, `true`: yes. Default value: `false`.
* `disable_higher_video_resolution` - (Optional, Bool) Whether to prohibit transcoding from low resolution to high resolution. Valid values: `false`,`true`. `false`: no, `true`: yes. Default value: `false`.
* `drm_type` - (Optional, String, ForceNew) DRM scheme type. Valid values: `SimpleAES`. If this field is an empty string, DRM will not be performed on the video.
* `segment_type` - (Optional, String) Segment type, valid when Format is HLS, optional values:
- ts: ts segment;
- fmp4: fmp4 segment;
Default value: ts.
* `sub_app_id` - (Optional, Int) The VOD [application](https://intl.cloud.tencent.com/document/product/266/14574) ID. For customers who activate VOD service from December 25, 2023, if they want to access resources in a VOD application (whether it's the default application or a newly created one), they must fill in this field with the application ID.

The `audio` object of `stream_info` supports the following:

* `bitrate` - (Required, Int) Audio stream bitrate in Kbps. Value range: `0` and `[26, 256]`. If the value is `0`, the bitrate of the audio stream will be the same as that of the original audio.
* `codec` - (Required, String) Audio stream encoder. Valid value are: `libfdk_aac` and `libmp3lame`. while `libfdk_aac` is recommended.
* `sample_rate` - (Required, Int) Audio stream sample rate. Valid values: `32000`, `44100`, `48000`Hz.
* `audio_channel` - (Optional, String) Audio channel system. Valid values: mono, dual, stereo. Default value: dual.

The `stream_info` object supports the following:

* `audio` - (Required, List) Audio parameter information.
* `video` - (Required, List) Video parameter information.
* `remove_audio` - (Optional, Bool) Whether to remove audio stream. Valid values: `false`: no, `true`: yes. `false` by default.
* `remove_video` - (Optional, Bool) Whether to remove video stream. Valid values: `false`: no, `true`: yes. `false` by default.
* `tehd_config` - (Optional, List) Extremely fast HD transcoding parameters.

The `tehd_config` object of `stream_info` supports the following:

* `type` - (Required, String) Extreme high-speed HD type, available values:
- TEHD-100: super high definition-100th;
- OFF: turn off Ultra High definition.
* `max_video_bitrate` - (Optional, Int) Video bitrate limit, which is valid when Type specifies extreme speed HD type. If you leave it empty or enter 0, there is no video bitrate limit.

The `video` object of `stream_info` supports the following:

* `bitrate` - (Required, Int) Bitrate of video stream in Kbps. Value range: `0` and `[128, 35000]`. If the value is `0`, the bitrate of the video will be the same as that of the source video.
* `codec` - (Required, String) Video stream encoder. Valid values: `libx264`,`libx265`,`av1`. `libx264`: H.264, `libx265`: H.265, `av1`: AOMedia Video 1. Currently, a resolution within 640x480 must be specified for `H.265`. and the `av1` container only supports mp4.
* `fps` - (Required, Int) Video frame rate in Hz. Value range: `[0, 60]`. If the value is `0`, the frame rate will be the same as that of the source video.
* `codec_tag` - (Optional, String) Encoding label, valid only if the encoding format of the video stream is H.265 encoding. Available values:
- hvc1: stands for hvc1 tag;
- hev1: stands for the hev1 tag;
Default value: hvc1.
* `fill_type` - (Optional, String) Fill type. Fill refers to the way of processing a screenshot when its aspect ratio is different from that of the source video. The following fill types are supported: `stretch`: stretch. The screenshot will be stretched frame by frame to match the aspect ratio of the source video, which may make the screenshot shorter or longer; `black`: fill with black. This option retains the aspect ratio of the source video for the screenshot and fills the unmatched area with black color blocks. Default value: black. Note: this field may return null, indicating that no valid values can be obtained.
* `gop` - (Optional, Int) Interval between Keyframe I frames, value range: 0 and [1, 100000], unit: number of frames. When you fill in 0 or leave it empty, the gop length is automatically set.
* `height` - (Optional, Int) Maximum value of the height (or short side) of a video stream in px. Value range: `0` and `[128, 4096]`. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Default value: `0`. Note: this field may return null, indicating that no valid values can be obtained.
* `preserve_hdr_switch` - (Optional, String) Whether the transcoding output still maintains HDR when the original video is HDR (High Dynamic Range). Value range:
- ON: if the original file is HDR, the transcoding output remains HDR;, otherwise the transcoding output is SDR (Standard Dynamic Range);
- OFF: regardless of whether the original file is HDR or SDR, the transcoding output is SDR;
Default value: OFF.
* `resolution_adaptive` - (Optional, Bool) Resolution adaption. Valid values: `true`,`false`. `true`: enabled. In this case, `width` represents the long side of a video, while `height` the short side; `false`: disabled. In this case, `width` represents the width of a video, while `height` the height. Default value: `true`. Note: this field may return null, indicating that no valid values can be obtained.
* `vcrf` - (Optional, Int) Video constant bit rate control factor, value range is [1,51].
Note:
- If this parameter is specified, the bitrate control method of CRF will be used for transcoding (the video bitrate will no longer take effect);
- This field is required when the video stream encoding format is H.266. The recommended value is 28;
- If there are no special requirements, it is not recommended to specify this parameter.
* `width` - (Optional, Int) Maximum value of the width (or long side) of a video stream in px. Value range: `0` and `[128, 4096]`. If both `width` and `height` are `0`, the resolution will be the same as that of the source video; If `width` is `0`, but `height` is not `0`, `width` will be proportionally scaled; If `width` is not `0`, but `height` is `0`, `height` will be proportionally scaled; If both `width` and `height` are not `0`, the custom resolution will be used. Default value: `0`. Note: this field may return null, indicating that no valid values can be obtained.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of template in ISO date format.
* `update_time` - Last modified time of template in ISO date format.


## Import

VOD adaptive dynamic streaming template can be imported using the id($subAppId#$templateId), e.g.

```
$ terraform import tencentcloud_vod_adaptive_dynamic_streaming_template.foo $subAppId#$templateId
```

