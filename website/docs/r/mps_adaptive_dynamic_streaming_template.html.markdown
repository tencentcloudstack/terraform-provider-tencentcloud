---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_adaptive_dynamic_streaming_template"
sidebar_current: "docs-tencentcloud-resource-mps_adaptive_dynamic_streaming_template"
description: |-
  Provides a resource to create a mps adaptive_dynamic_streaming_template
---

# tencentcloud_mps_adaptive_dynamic_streaming_template

Provides a resource to create a mps adaptive_dynamic_streaming_template

## Example Usage

```hcl
resource "tencentcloud_mps_adaptive_dynamic_streaming_template" "adaptive_dynamic_streaming_template" {
  comment                         = "terrraform test"
  disable_higher_video_bitrate    = 0
  disable_higher_video_resolution = 1
  format                          = "HLS"
  name                            = "terrraform-test"

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

## Argument Reference

The following arguments are supported:

* `format` - (Required, String) Adaptive transcoding format, value range:HLS, MPEG-DASH.
* `stream_infos` - (Required, List) Convert adaptive code stream to output sub-stream parameter information, and output up to 10 sub-streams.Note: The frame rate of each sub-stream must be consistent; if not, the frame rate of the first sub-stream is used as the output frame rate.
* `comment` - (Optional, String) Template description information, length limit: 256 characters.
* `disable_higher_video_bitrate` - (Optional, Int) Whether to prohibit video from low bit rate to high bit rate, value range:0: no.1: yes.Default value: 0.
* `disable_higher_video_resolution` - (Optional, Int) Whether to prohibit the conversion of video resolution to high resolution, value range:0: no.1: yes.Default value: 0.
* `name` - (Optional, String) Template name, length limit: 64 characters.

The `audio` object supports the following:

* `bitrate` - (Required, Int) Bit rate of the audio stream, value range: 0 and [26, 256], unit: kbps.When the value is 0, it means that the audio bit rate is consistent with the original audio.
* `codec` - (Required, String) Encoding format of audio stream.When the outer parameter Container is mp3, the optional value is:libmp3lame.When the outer parameter Container is ogg or flac, the optional value is:flac.When the outer parameter Container is m4a, the optional value is:libfdk_aac.libmp3lame.ac3.When the outer parameter Container is mp4 or flv, the optional value is:libfdk_aac: more suitable for mp4.libmp3lame: more suitable for flv.When the outer parameter Container is hls, the optional value is:libfdk_aac.libmp3lame.
* `sample_rate` - (Required, Int) Sampling rate of audio stream, optional value.32000.44100.48000.Unit: Hz.
* `audio_channel` - (Optional, Int) Audio channel mode, optional values:`1: single channel.2: Dual channel.6: Stereo.When the package format of the media is an audio format (flac, ogg, mp3, m4a), the number of channels is not allowed to be set to stereo.Default: 2.

The `stream_infos` object supports the following:

* `audio` - (Required, List) Audio parameter information.
* `video` - (Required, List) Video parameter information.
* `remove_audio` - (Optional, Int) Whether to remove audio stream, value:0: reserved.1: remove.
* `remove_video` - (Optional, Int) Whether to remove video stream, value:0: reserved.1: remove.

The `video` object supports the following:

* `bitrate` - (Required, Int) Bit rate of the video stream, value range: 0 and [128, 35000], unit: kbps.When the value is 0, it means that the video bit rate is consistent with the original video.
* `codec` - (Required, String) Encoding format of the video stream, optional value:libx264: H.264 encoding.libx265: H.265 encoding.av1: AOMedia Video 1 encoding.Note: Currently H.265 encoding must specify a resolution, and it needs to be within 640*480.Note: av1 encoded containers currently only support mp4.
* `fps` - (Required, Int) Video frame rate, value range: [0, 100], unit: Hz.When the value is 0, it means that the frame rate is consistent with the original video.Note: The value range for adaptive code rate is [0, 60].
* `fill_type` - (Optional, String) Filling type, when the aspect ratio of the video stream configuration is inconsistent with the aspect ratio of the original video, the processing method for transcoding is filling. Optional filling type:stretch: Stretching, stretching each frame to fill the entire screen, which may cause the transcoded video to be squashed or stretched.black: Leave black, keep the video aspect ratio unchanged, and fill the rest of the edge with black.white: Leave blank, keep the aspect ratio of the video, and fill the rest of the edge with white.gauss: Gaussian blur, keep the aspect ratio of the video unchanged, and use Gaussian blur for the rest of the edge.Default value: black.Note: Adaptive stream only supports stretch, black.
* `gop` - (Optional, Int) The interval between keyframe I frames, value range: 0 and [1, 100000], unit: number of frames.When filling 0 or not filling, the system will automatically set the gop length.
* `height` - (Optional, Int) The maximum value of the height (or short side) of the video streaming, value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.
* `resolution_adaptive` - (Optional, String) Adaptive resolution, optional value:open: At this time, Width represents the long side of the video, Height represents the short side of the video.close: At this point, Width represents the width of the video, and Height represents the height of the video.Default value: open.Note: In adaptive mode, Width cannot be smaller than Height.
* `vcrf` - (Optional, Int) Video constant bit rate control factor, the value range is [1, 51].If this parameter is specified, the code rate control method of CRF will be used for transcoding (the video code rate will no longer take effect).If there is no special requirement, it is not recommended to specify this parameter.
* `width` - (Optional, Int) The maximum value of the width (or long side) of the video streaming, value range: 0 and [128, 4096], unit: px.When Width and Height are both 0, the resolution is the same.When Width is 0 and Height is not 0, Width is scaled proportionally.When Width is not 0 and Height is 0, Height is scaled proportionally.When both Width and Height are not 0, the resolution is specified by the user.Default value: 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps adaptive_dynamic_streaming_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_adaptive_dynamic_streaming_template.adaptive_dynamic_streaming_template adaptive_dynamic_streaming_template_id
```

