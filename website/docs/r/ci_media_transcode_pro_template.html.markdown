---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_media_transcode_pro_template"
sidebar_current: "docs-tencentcloud-resource-ci_media_transcode_pro_template"
description: |-
  Provides a resource to create a ci media_transcode_pro_template
---

# tencentcloud_ci_media_transcode_pro_template

Provides a resource to create a ci media_transcode_pro_template

## Example Usage

```hcl
resource "tencentcloud_ci_media_transcode_pro_template" "media_transcode_pro_template" {
  bucket = "terraform-ci-xxxxxx"
  name   = "transcode_pro_template"
  container {
    format = "mxf"
    # clip_config {
    # 	duration = ""
    # }

  }
  video {
    codec      = "xavc"
    profile    = "XAVC-HD_422_10bit"
    width      = "1920"
    height     = "1080"
    interlaced = "true"
    fps        = "30000/1001"
    bitrate    = "50000"
    # rotate = ""

  }
  time_interval {
    start    = ""
    duration = ""

  }
  audio {
    codec  = "pcm_s24le"
    remove = "true"

  }
  trans_config {
    adj_dar_method           = "scale"
    is_check_reso            = "false"
    reso_adj_method          = "1"
    is_check_video_bitrate   = "false"
    video_bitrate_adj_method = "0"
    is_check_audio_bitrate   = "false"
    audio_bitrate_adj_method = "0"
    delete_metadata          = "false"
    is_hdr2_sdr              = "false"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) bucket name.
* `container` - (Required, List) container format.
* `name` - (Required, String) The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.
* `audio` - (Optional, List) Audio information, do not transmit Audio, which is equivalent to deleting audio information.
* `time_interval` - (Optional, List) time interval.
* `trans_config` - (Optional, List) transcoding configuration.
* `video` - (Optional, List) video information, do not upload Video, which is equivalent to deleting video information.

The `audio` object supports the following:

* `codec` - (Required, String) Codec format, value aac, mp3, flac, amr, Vorbis, opus, pcm_s16le.
* `remove` - (Optional, String) Whether to delete the source audio stream, the value is true, false.

The `clip_config` object supports the following:

* `duration` - (Optional, String) Fragmentation duration, default 5s.

The `container` object supports the following:

* `format` - (Required, String) Package format.
* `clip_config` - (Optional, List) Fragment configuration, valid when format is hls and dash.

The `time_interval` object supports the following:

* `duration` - (Optional, String) duration, [0 video duration], in seconds, Support float format, the execution accuracy is accurate to milliseconds.
* `start` - (Optional, String) Starting time, [0 video duration], in seconds, Support float format, the execution accuracy is accurate to milliseconds.

The `trans_config` object supports the following:

* `adj_dar_method` - (Optional, String) Resolution adjustment method, value scale, crop, pad, none, When the aspect ratio of the output video is different from the original video, adjust the resolution accordingly according to this parameter.
* `audio_bitrate_adj_method` - (Optional, String) Audio bit rate adjustment mode, value 0, 1; when the output audio bit rate is greater than the original audio bit rate, 0 means use the original audio bit rate; 1 means return transcoding failed, Take effect when IsCheckAudioBitrate is true.
* `delete_metadata` - (Optional, String) Whether to delete the MetaData information in the file, true, false, When false, keep source file information.
* `is_check_audio_bitrate` - (Optional, String) Whether to check the audio code rate, true, false, When false, transcode according to configuration parameters.
* `is_check_reso` - (Optional, String) Whether to check the resolution, when it is false, transcode according to the configuration parameters.
* `is_check_video_bitrate` - (Optional, String) Whether to check the video code rate, when it is false, transcode according to the configuration parameters.
* `is_hdr2_sdr` - (Optional, String) Whether to enable HDR to SDR true, false.
* `reso_adj_method` - (Optional, String) Resolution adjustment mode, value 0, 1; 0 means use the original video resolution; 1 means return transcoding failed, Take effect when IsCheckReso is true.
* `video_bitrate_adj_method` - (Optional, String) Video bit rate adjustment method, value 0, 1; when the output video bit rate is greater than the original video bit rate, 0 means use the original video bit rate; 1 means return transcoding failed, Take effect when IsCheckVideoBitrate is true.

The `video` object supports the following:

* `bitrate` - (Optional, String) Bit rate of video output file, value range: [10, 50000], unit: Kbps, auto means adaptive bit rate.
* `codec` - (Optional, String) Codec format, default value: `H.264`, when format is WebM, it is VP8, value range: `H.264`, `H.265`, `VP8`, `VP9`, `AV1`.
* `fps` - (Optional, String) Frame rate, value range: (0, 60], Unit: fps.
* `height` - (Optional, String) High, value range: [128, 4096], Unit: px, If only Height is set, Width is calculated according to the original ratio of the video, must be even.
* `interlaced` - (Optional, String) field pattern.
* `profile` - (Optional, String) encoding level, Support baseline, main, high, auto- When Pixfmt is auto, this parameter can only be set to auto, when it is set to other options, the parameter value will be set to auto- baseline: suitable for mobile devices- main: suitable for standard resolution devices- high: suitable for high-resolution devices- Only H.264 supports this parameter.
* `rotate` - (Optional, String) Rotation angle, Value range: [0, 360), Unit: degree.
* `width` - (Optional, String) width, value range: [128, 4096], Unit: px, If only Width is set, Height is calculated according to the original ratio of the video, must be even.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci media_transcode_pro_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_transcode_pro_template.media_transcode_pro_template terraform-ci-xxxxxx#t13ed9af009da0414e9c7c63456ec8f4d2
```

