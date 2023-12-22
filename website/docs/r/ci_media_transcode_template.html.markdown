---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_media_transcode_template"
sidebar_current: "docs-tencentcloud-resource-ci_media_transcode_template"
description: |-
  Provides a resource to create a ci media_transcode_template
---

# tencentcloud_ci_media_transcode_template

Provides a resource to create a ci media_transcode_template

## Example Usage

```hcl
resource "tencentcloud_ci_media_transcode_template" "media_transcode_template" {
  bucket = "terraform-ci-1308919341"
  name   = "transcode_template"
  container {
    format = "mp4"
    # clip_config {
    # 	duration = ""
    # }
  }
  video {
    codec = "H.264"
    width = "1280"
    # height = ""
    fps     = "30"
    remove  = "false"
    profile = "high"
    bitrate = "1000"
    # crf = ""
    # gop = ""
    preset = "medium"
    # bufsize = ""
    # maxrate = ""
    # pixfmt = ""
    long_short_mode = "false"
    # rotate = ""
  }
  time_interval {
    start    = "0"
    duration = "60"
  }
  audio {
    codec           = "aac"
    samplerate      = "44100"
    bitrate         = "128"
    channels        = "4"
    remove          = "false"
    keep_two_tracks = "false"
    switch_track    = "false"
    sample_format   = ""
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
  audio_mix {
    audio_source = "https://terraform-ci-1308919341.cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"
    mix_mode     = "Once"
    replace      = "true"
    effect_config {
      enable_start_fadein = "true"
      start_fadein_time   = "3"
      enable_end_fadeout  = "false"
      end_fadeout_time    = "0"
      enable_bgm_fade     = "true"
      bgm_fade_time       = "1.7"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) bucket name.
* `container` - (Required, List) container format.
* `name` - (Required, String) The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.
* `audio_mix` - (Optional, List) mixing parameters.
* `audio` - (Optional, List) Audio information, do not transmit Audio, which is equivalent to deleting audio information.
* `time_interval` - (Optional, List) time interval.
* `trans_config` - (Optional, List) transcoding configuration.
* `video` - (Optional, List) video information, do not upload Video, which is equivalent to deleting video information.

The `audio_mix` object supports the following:

* `audio_source` - (Required, String) The media address of the audio track that needs to be mixed.
* `effect_config` - (Optional, List) Mix Fade Configuration.
* `mix_mode` - (Optional, String) Mixing mode Repeat: background sound loop, Once: The background sound is played once.
* `replace` - (Optional, String) Whether to replace the original audio of the Input media file with the mixed audio track media.

The `audio` object supports the following:

* `bitrate` - (Optional, String) Original audio bit rate, unit: Kbps, Value range: [8, 1000].
* `channels` - (Optional, String) number of channels- When Codec is set to aac/flac, support 1, 2, 4, 5, 6, 8- When Codec is set to mp3/opus, support 1, 2- When Codec is set to Vorbis, only 2 is supported- When Codec is set to amr, only 1 is supported- When Codec is set to pcm_s16le, only 1 and 2 are supported- When the encapsulation format is dash, 8 is not supported.
* `codec` - (Optional, String) Codec format, value aac, mp3, flac, amr, Vorbis, opus, pcm_s16le.
* `keep_two_tracks` - (Optional, String) Keep dual audio tracks, the value is true, false. This parameter is invalid when Video.Codec is H.265.
* `remove` - (Optional, String) Whether to delete the source audio stream, the value is true, false.
* `sample_format` - (Optional, String) Sampling bit width- When Codec is set to aac, support fltp- When Codec is set to mp3, fltp, s16p, s32p are supported- When Codec is set to flac, s16, s32, s16p, s32p are supported- When Codec is set to amr, support s16, s16p- When Codec is set to opus, support s16- When Codec is set to pcm_s16le, support s16- When Codec is set to Vorbis, support fltp- This parameter is invalid when Video.Codec is H.265.
* `samplerate` - (Optional, String) Sampling Rate- Unit: Hz- Optional 8000, 11025, 12000, 16000, 22050, 24000, 32000, 44100, 48000, 88200, 96000- Different packages, mp3 supports different sampling rates, as shown in the table below- When Codec is set to amr, only 8000 is supported- When Codec is set to opus, it supports 8000, 16000, 24000, 48000.
* `switch_track` - (Optional, String) Convert track, the value is true, false. This parameter is invalid when Video.Codec is H.265.

The `clip_config` object of `container` supports the following:

* `duration` - (Optional, String) Fragmentation duration, default 5s.

The `container` object supports the following:

* `format` - (Required, String) Package format.
* `clip_config` - (Optional, List) Fragment configuration, valid when format is hls and dash.

The `effect_config` object of `audio_mix` supports the following:

* `bgm_fade_time` - (Optional, String) bgm transition fade-in duration, support floating point numbers.
* `enable_bgm_fade` - (Optional, String) Enable bgm conversion fade in.
* `enable_end_fadeout` - (Optional, String) enable fade out.
* `enable_start_fadein` - (Optional, String) enable fade in.
* `end_fadeout_time` - (Optional, String) fade out time, greater than 0, support floating point numbers.
* `start_fadein_time` - (Optional, String) Fade in duration, greater than 0, support floating point numbers.

The `hls_encrypt` object of `trans_config` supports the following:

* `is_hls_encrypt` - (Optional, String) Whether to enable HLS encryption, support encryption when Container.Format is hls.
* `uri_key` - (Optional, String) HLS encrypted key, this parameter is only meaningful when IsHlsEncrypt is true.

The `time_interval` object supports the following:

* `duration` - (Optional, String) duration, [0 video duration], in seconds, Support float format, the execution accuracy is accurate to milliseconds.
* `start` - (Optional, String) Starting time, [0 video duration], in seconds, Support float format, the execution accuracy is accurate to milliseconds.

The `trans_config` object supports the following:

* `adj_dar_method` - (Optional, String) Resolution adjustment method, value scale, crop, pad, none, When the aspect ratio of the output video is different from the original video, adjust the resolution accordingly according to this parameter.
* `audio_bitrate_adj_method` - (Optional, String) Audio bit rate adjustment mode, value 0, 1; when the output audio bit rate is greater than the original audio bit rate, 0 means use the original audio bit rate; 1 means return transcoding failed, Take effect when IsCheckAudioBitrate is true.
* `delete_metadata` - (Optional, String) Whether to delete the MetaData information in the file, true, false, When false, keep source file information.
* `hls_encrypt` - (Optional, List) hls encryption configuration.
* `is_check_audio_bitrate` - (Optional, String) Whether to check the audio code rate, true, false, When false, transcode according to configuration parameters.
* `is_check_reso` - (Optional, String) Whether to check the resolution, when it is false, transcode according to the configuration parameters.
* `is_check_video_bitrate` - (Optional, String) Whether to check the video code rate, when it is false, transcode according to the configuration parameters.
* `is_hdr2_sdr` - (Optional, String) Whether to enable HDR to SDR true, false.
* `reso_adj_method` - (Optional, String) Resolution adjustment mode, value 0, 1; 0 means use the original video resolution; 1 means return transcoding failed, Take effect when IsCheckReso is true.
* `video_bitrate_adj_method` - (Optional, String) Video bit rate adjustment method, value 0, 1; when the output video bit rate is greater than the original video bit rate, 0 means use the original video bit rate; 1 means return transcoding failed, Take effect when IsCheckVideoBitrate is true.

The `video` object supports the following:

* `bitrate` - (Optional, String) Bit rate of video output file, value range: [10, 50000], unit: Kbps, auto means adaptive bit rate.
* `bufsize` - (Optional, String) buffer size, Value range: [1000, 128000], Unit: Kb, This parameter is not supported when Codec is VP8/VP9.
* `codec` - (Optional, String) Codec format, default value: `H.264`, when format is WebM, it is VP8, value range: `H.264`, `H.265`, `VP8`, `VP9`, `AV1`.
* `crf` - (Optional, String) Bit rate-quality control factor, value range: (0, 51], If Crf is set, the setting of Bitrate will be invalid, When Bitrate is empty, the default is 25.
* `fps` - (Optional, String) Frame rate, value range: (0, 60], Unit: fps.
* `gop` - (Optional, String) The maximum number of frames between key frames, value range: [1, 100000].
* `height` - (Optional, String) High, value range: [128, 4096], Unit: px, If only Height is set, Width is calculated according to the original ratio of the video, must be even.
* `long_short_mode` - (Optional, String) Adaptive length,true, false, This parameter is not supported when Codec is VP8/VP9/AV1.
* `maxrate` - (Optional, String) Peak video bit rate, Value range: [10, 50000], Unit: Kbps, This parameter is not supported when Codec is VP8/VP9.
* `pixfmt` - (Optional, String) video color format, H.264 support: yuv420p, yuv422p, yuv444p, yuvj420p, yuvj422p, yuvj444p, auto, H.265 support: yuv420p, yuv420p10le, auto, This parameter is not supported when Codec is VP8/VP9/AV1.
* `preset` - (Optional, String) Video Algorithm Presets- H.264 supports this parameter, the values are veryfast, fast, medium, slow, slower- VP8 supports this parameter, the value is good, realtime- AV1 supports this parameter, the value is 5 (recommended value), 4- H.265 and VP9 do not support this parameter.
* `profile` - (Optional, String) encoding level, Support baseline, main, high, auto- When Pixfmt is auto, this parameter can only be set to auto, when it is set to other options, the parameter value will be set to auto- baseline: suitable for mobile devices- main: suitable for standard resolution devices- high: suitable for high-resolution devices- Only H.264 supports this parameter.
* `remove` - (Optional, String) Whether to delete the video stream, true, false.
* `rotate` - (Optional, String) Rotation angle, Value range: [0, 360), Unit: degree.
* `width` - (Optional, String) width, value range: [128, 4096], Unit: px, If only Width is set, Height is calculated according to the original ratio of the video, must be even.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci media_transcode_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_transcode_template.media_transcode_template media_transcode_template_id
```

