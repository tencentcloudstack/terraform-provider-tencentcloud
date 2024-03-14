---
subcategory: "Video on Demand(VOD)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_transcode_template"
sidebar_current: "docs-tencentcloud-resource-vod_transcode_template"
description: |-
  Provides a resource to create a vod transcode template
---

# tencentcloud_vod_transcode_template

Provides a resource to create a vod transcode template

## Example Usage

```hcl
resource "tencentcloud_vod_sub_application" "sub_application" {
  name        = "transcodeTemplateSubApplication"
  status      = "On"
  description = "this is sub application"
}

resource "tencentcloud_vod_transcode_template" "transcode_template" {
  container    = "mp4"
  sub_app_id   = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  name         = "720pTranscodeTemplate"
  comment      = "test transcode mp4 720p"
  remove_video = 0
  remove_audio = 0
  video_template {
    codec               = "libx264"
    fps                 = 26
    bitrate             = 1000
    resolution_adaptive = "open"
    width               = 0
    height              = 720
    fill_type           = "stretch"
    vcrf                = 1
    gop                 = 250
    preserve_hdr_switch = "OFF"
    codec_tag           = "hvc1"

  }
  audio_template {
    codec         = "libfdk_aac"
    bitrate       = 128
    sample_rate   = 44100
    audio_channel = 2
  }
  segment_type = "ts"
}
```

## Argument Reference

The following arguments are supported:

* `container` - (Required, String) The container format. Valid values: `mp4`, `flv`, `hls`, `mp3`, `flac`, `ogg`, `m4a`, `wav` ( `mp3`, `flac`, `ogg`, `m4a`, and `wav` are audio file formats).
* `audio_template` - (Optional, List) Audio stream configuration parameter. This field is required when `RemoveAudio` is 0.
* `comment` - (Optional, String) Template description. Length limit: 256 characters.
* `name` - (Optional, String) Transcoding template name. Length limit: 64 characters.
* `remove_audio` - (Optional, Int) Whether to remove audio data. Valid values:0: retain 1: remove Default value: 0.
* `remove_video` - (Optional, Int) Whether to remove video data. Valid values:
- 0: retain
- 1: remove
Default value: 0.
* `segment_type` - (Optional, String) The segment type. This parameter is valid only if `Container` is `hls`. Valid values: `ts`: TS segment; `fmp4`: fMP4 segment Default: `ts`.
* `sub_app_id` - (Optional, Int) The VOD [application](https://intl.cloud.tencent.com/document/product/266/14574) ID. For customers who activate VOD service from December 25, 2023, if they want to access resources in a VOD application (whether it's the default application or a newly created one), they must fill in this field with the application ID.
* `tehd_config` - (Optional, List) TESHD transcoding parameter.
* `video_template` - (Optional, List) Video stream configuration parameter. This field is required when `RemoveVideo` is 0.

The `audio_template` object supports the following:

* `bitrate` - (Required, Int) Audio stream bitrate in Kbps. Value range: 0 and [26, 256].If the value is 0, the bitrate of the audio stream will be the same as that of the original audio.
* `codec` - (Required, String) The audio codec.If `Container` is `mp3`, the valid value is:`libmp3lame`If `Container` is `ogg` or `flac`, the valid value is:`flac`If `Container` is `m4a`, the valid values are:`libfdk_aac``libmp3lame``ac3`If `Container` is `mp4` or `flv`, the valid values are:`libfdk_aac` (Recommended for MP4)`libmp3lame` (Recommended for FLV)`mp2`If `Container` is `hls`, the valid value is:`libfdk_aac`If `Format` is `HLS` or `MPEG-DASH`, the valid value is:`libfdk_aac`If `Container` is `wav`, the valid value is:`pcm16`.
* `sample_rate` - (Required, Int) The audio sample rate. Valid values:`16000` (valid only if `Codec` is `pcm16`)`32000``44100``48000`Unit: Hz.
* `audio_channel` - (Optional, Int) Audio channel system. Valid values:1: mono-channel2: dual-channel6: stereoYou cannot set the sound channel as stereo for media files in container formats for audios (FLAC, OGG, MP3, M4A).Default value: 2.

The `tehd_config` object supports the following:

* `type` - (Required, String) TESHD transcoding type. Valid values: TEHD-100, OFF (default).
* `max_video_bitrate` - (Optional, Int) Maximum bitrate, which is valid when `Type` is `TESHD`.If this parameter is left blank or 0 is entered, there will be no upper limit for bitrate.

The `video_template` object supports the following:

* `bitrate` - (Required, Int) Bitrate of video stream in Kbps. Value range: 0 and [128, 35,000].If the value is 0, the bitrate of the video will be the same as that of the source video.
* `codec` - (Required, String) The video codec. Valid values:libx264: H.264; libx265: H.265; av1: AOMedia Video 1; H.266: H.266. The AOMedia Video 1 and H.266 codecs can only be used for MP4 files. Only CRF is supported for H.266 currently.
* `fps` - (Required, Int) Video frame rate in Hz. Value range: [0,100].If the value is 0, the frame rate will be the same as that of the source video.
* `codec_tag` - (Optional, String) The codec tag. This parameter is valid only if the H.265 codec is used. Valid values:hvc1hev1Default value: hvc1.
* `fill_type` - (Optional, String) Fill type, the way of processing a screenshot when the configured aspect ratio is different from that of the source video. Valid values:stretch: stretches the video image frame by frame to fill the screen. The video image may become squashed or stretched after transcoding.black: fills the uncovered area with black color, without changing the image&#39;s aspect ratio.white: fills the uncovered area with white color, without changing the image&#39;s aspect ratio.gauss: applies Gaussian blur to the uncovered area, without changing the image&#39;s aspect ratio.Default value: black.
* `gop` - (Optional, Int) I-frame interval in frames. Valid values: 0 and 1-100000.When this parameter is set to 0 or left empty, `Gop` will be automatically set.
* `height` - (Optional, Int) The maximum video height (or short side) in pixels. Value range: 0 and [128, 8192].If both `Width` and `Height` are 0, the output resolution will be the same as that of the source video.If `Width` is 0 and `Height` is not, the video width will be proportionally scaled.If `Width` is not 0 and `Height` is, the video height will be proportionally scaled.If neither `Width` nor `Height` is 0, the specified width and height will be used.Default value: 0.
* `preserve_hdr_switch` - (Optional, String) Whether to output an HDR (high dynamic range) video if the source video is HDR. Valid values:ON: If the source video is HDR, output an HDR video; if not, output an SDR (standard dynamic range) video.OFF: Output an SDR video regardless of whether the source video is HDR.Default value: OFF.
* `resolution_adaptive` - (Optional, String) Resolution adaption. Valid values:open: enabled. In this case, `Width` represents the long side of a video, while `Height` the short side;close: disabled. In this case, `Width` represents the width of a video, while `Height` the height.Default value: open.Note: this field may return null, indicating that no valid values can be obtained.
* `vcrf` - (Optional, Int) The video constant rate factor (CRF). Value range: 1-51.If this parameter is specified, CRF encoding will be used and the bitrate parameter will be ignored.If `Codec` is `H.266`, this parameter is required (`28` is recommended).We don't recommend using this parameter unless you have special requirements.
* `width` - (Optional, Int) The maximum video width (or long side) in pixels. Value range: 0 and [128, 8192].If both `Width` and `Height` are 0, the output resolution will be the same as that of the source video.If `Width` is 0 and `Height` is not, the video width will be proportionally scaled.If `Width` is not 0 and `Height` is, the video height will be proportionally scaled.If neither `Width` nor `Height` is 0, the specified width and height will be used.Default value: 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

vod transcode template can be imported using the id, e.g.

```
terraform import tencentcloud_vod_transcode_template.transcode_template $subAppId#$templateId
```

