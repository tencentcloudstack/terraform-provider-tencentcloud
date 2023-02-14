---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_transcode_template"
sidebar_current: "docs-tencentcloud-resource-mps_transcode_template"
description: |-
  Provides a resource to create a mps transcode_template
---

# tencentcloud_mps_transcode_template

Provides a resource to create a mps transcode_template

## Example Usage

```hcl
resource "tencentcloud_mps_transcode_template" "transcode_template" {
  container    = "mp4"
  name         = "tf_transcode_template"
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

## Argument Reference

The following arguments are supported:

* `container` - (Required, String) Encapsulation format, optional values: mp4, flv, hls, mp3, flac, ogg, m4a. Among them, mp3, flac, ogg, m4a are pure audio files.
* `audio_template` - (Optional, List) Audio stream configuration parameters, when RemoveAudio is 0, this field is required.
* `comment` - (Optional, String) Template description information, length limit: 256 characters.
* `enhance_config` - (Optional, List) Audio and video enhancement configuration.
* `name` - (Optional, String) Transcoding template name, length limit: 64 characters.
* `remove_audio` - (Optional, Int) Whether to remove audio data, value:0: reserved.1: remove.Default: 0.
* `remove_video` - (Optional, Int) Whether to remove video data, value:0: reserved.1: remove.Default: 0.
* `tehd_config` - (Optional, List) Ultra-fast HD transcoding parameters.
* `video_template` - (Optional, List) Video stream configuration parameters, when RemoveVideo is 0, this field is required.

The `artifact_repair` object supports the following:

* `switch` - (Optional, String) Capability configuration switch, optional value: ON/OFF.Default value: ON.
* `type` - (Optional, String) Type, optional value: weak/strong.Default value: weak.Note: This field may return null, indicating that no valid value can be obtained.

The `audio_template` object supports the following:

* `bitrate` - (Required, Int) Bit rate of the audio stream, value range: 0 and [26, 256], unit: kbps.When the value is 0, it means that the audio bit rate is consistent with the original audio.
* `codec` - (Required, String) Encoding format of frequency stream.When the outer parameter Container is mp3, the optional value is:libmp3lame.When the outer parameter Container is ogg or flac, the optional value is:flac.When the outer parameter Container is m4a, the optional value is:libfdk_aac.libmp3lame.ac3.When the outer parameter Container is mp4 or flv, the optional value is:libfdk_aac: more suitable for mp4.libmp3lame: more suitable for flv.When the outer parameter Container is hls, the optional value is:libfdk_aac.libmp3lame.
* `sample_rate` - (Required, Int) Sampling rate of audio stream, optional value.32000.44100.48000.Unit: Hz.
* `audio_channel` - (Optional, Int) Audio channel mode, optional values:`1: single channel.2: Dual channel.6: Stereo.When the package format of the media is an audio format (flac, ogg, mp3, m4a), the number of channels is not allowed to be set to stereo.Default: 2.

The `color_enhance` object supports the following:

* `switch` - (Optional, String) Capability configuration switch, optional value: ON/OFF.Default value: ON.
* `type` - (Optional, String) Type, optional value: weak/normal/strong.Default value: weak.Note: This field may return null, indicating that no valid value can be obtained.

The `denoise` object supports the following:

* `switch` - (Optional, String) Capability configuration switch, optional value: ON/OFF.Default value: ON.
* `type` - (Optional, String) Type, optional value: weak/strong.Default value: weak.Note: This field may return null, indicating that no valid value can be obtained.

The `enhance_config` object supports the following:

* `video_enhance` - (Optional, List) Video Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.

The `face_enhance` object supports the following:

* `intensity` - (Optional, Float64) Intensity, value range: 0.0~1.0.Default value: 0.0.Note: This field may return null, indicating that no valid value can be obtained.
* `switch` - (Optional, String) Capability configuration switch, optional value: ON/OFF.Default value: ON.

The `frame_rate` object supports the following:

* `fps` - (Optional, Int) Frame rate, value range: [0, 100], unit: Hz.Default value: 0.Note: For transcoding, this parameter will override the Fps inside the VideoTemplate.Note: This field may return null, indicating that no valid value can be obtained.
* `switch` - (Optional, String) Capability configuration switch, optional value: ON/OFF.Default value: ON.

The `hdr` object supports the following:

* `switch` - (Optional, String) Capability configuration switch, optional value: ON/OFF.Default value: ON.
* `type` - (Optional, String) Type, optional value: HDR10/HLG.Default value: HDR10.Note: The encoding method of video needs to be libx265.Note: Video encoding bit depth is 10.Note: This field may return null, indicating that no valid value can be obtained.

The `image_quality_enhance` object supports the following:

* `switch` - (Optional, String) Capability configuration switch, optional value: ON/OFF.Default value: ON.
* `type` - (Optional, String) Type, optional value: weak/normal/strong.Default value: weak.Note: This field may return null, indicating that no valid value can be obtained.

The `low_light_enhance` object supports the following:

* `switch` - (Optional, String) Capability configuration switch, optional value: ON/OFF.Default value: ON.
* `type` - (Optional, String) Type, optional value: normal.Default value: normal.Note: This field may return null, indicating that no valid value can be obtained.

The `scratch_repair` object supports the following:

* `intensity` - (Optional, Float64) Intensity, value range: 0.0~1.0.Default value: 0.0.Note: This field may return null, indicating that no valid value can be obtained.
* `switch` - (Optional, String) Capability configuration switch, optional value: ON/OFF.Default value: ON.

The `sharp_enhance` object supports the following:

* `intensity` - (Optional, Float64) Intensity, value range: 0.0~1.0.Default value: 0.0.Note: This field may return null, indicating that no valid value can be obtained.
* `switch` - (Optional, String) Capability configuration switch, optional value: ON/OFF.Default value: ON.

The `super_resolution` object supports the following:

* `size` - (Optional, Int) Super resolution multiple, optional value:2: currently only supports 2x super resolution.Default value: 2.Note: This field may return null, indicating that no valid value can be obtained.
* `switch` - (Optional, String) Capability configuration switch, optional value: ON/OFF.Default value: ON.
* `type` - (Optional, String) Type, optional value:lq: super-resolution for low-definition video with more noise.hq: super resolution for high-definition video.Default value: lq.Note: This field may return null, indicating that no valid value can be obtained.

The `tehd_config` object supports the following:

* `type` - (Required, String) Extremely high-definition type, optional value:TEHD-100: Extreme HD-100.Not filling means that the ultra-fast high-definition is not enabled.
* `max_video_bitrate` - (Optional, Int) The upper limit of the video bit rate, which is valid when the Type specifies the ultra-fast HD type.Do not fill in or fill in 0 means that there is no upper limit on the video bit rate.

The `video_enhance` object supports the following:

* `artifact_repair` - (Optional, List) De-artifact (glitch) configuration.Note: This field may return null, indicating that no valid value can be obtained.
* `color_enhance` - (Optional, List) Color Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.
* `denoise` - (Optional, List) Video Noise Reduction Configuration.Note: This field may return null, indicating that no valid value can be obtained.
* `face_enhance` - (Optional, List) Face Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.
* `frame_rate` - (Optional, List) Interpolation frame rate configuration.Note: This field may return null, indicating that no valid value can be obtained.
* `hdr` - (Optional, List) HDR configuration.Note: This field may return null, indicating that no valid value can be obtained.
* `image_quality_enhance` - (Optional, List) Comprehensive Enhanced Configuration.Note: This field may return null, indicating that no valid value can be obtained.
* `low_light_enhance` - (Optional, List) Low Light Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.
* `scratch_repair` - (Optional, List) De-scratch configuration.Note: This field may return null, indicating that no valid value can be obtained.
* `sharp_enhance` - (Optional, List) Detail Enhancement Configuration.Note: This field may return null, indicating that no valid value can be obtained.
* `super_resolution` - (Optional, List) Super resolution configuration.Note: This field may return null, indicating that no valid value can be obtained.

The `video_template` object supports the following:

* `bitrate` - (Required, Int) Bit rate of the video stream, value range: 0 and [128, 35000], unit: kbps.When the value is 0, it means that the video bit rate is consistent with the original video.
* `codec` - (Required, String) Encoding format of the video stream, optional value:libx264: H.264 encoding.libx265: H.265 encoding.av1: AOMedia Video 1 encoding.Note: Currently H.265 encoding must specify a resolution, and it needs to be within 640*480.Note: av1 encoded containers currently only support mp4.
* `fps` - (Required, Int) Video frame rate, value range: [0, 100], unit: Hz.When the value is 0, it means that the frame rate is consistent with the original video.Note: The value range for adaptive code rate is [0, 60].
* `fill_type` - (Optional, String) Filling method, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling method:stretch: Stretch, stretch each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched.black: Leave black, keep the aspect ratio of the video unchanged, and fill the rest of the edge with black.white: Leave blank, keep the aspect ratio of the video unchanged, and fill the rest of the edge with white.gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and fill the rest of the edge with Gaussian blur.Default: black.Note: Adaptive stream only supports stretch, black.
* `gop` - (Optional, Int) The interval between keyframe I frames, value range: 0 and [1, 100000], unit: number of frames.When filling 0 or not filling, the system will automatically set the gop length.
* `height` - (Optional, Int) The maximum value of video stream height (or short side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default: 0.
* `resolution_adaptive` - (Optional, String) Adaptive resolution, optional values:```open: open, at this time, Width represents the long side of the video, Height represents the short side of the video.close: close, at this time, Width represents the width of the video, and Height represents the height of the video.Default: open.Note: In adaptive mode, Width cannot be smaller than Height.
* `vcrf` - (Optional, Int) Video constant bit rate control factor, the value range is [1, 51].If this parameter is specified, the code rate control method of CRF will be used for transcoding (the video code rate will no longer take effect).If there is no special requirement, it is not recommended to specify this parameter.
* `width` - (Optional, Int) The maximum value of video stream width (or long side), value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default: 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps transcode_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_transcode_template.transcode_template transcode_template_id
```

