---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_media_concat_template"
sidebar_current: "docs-tencentcloud-resource-ci_media_concat_template"
description: |-
  Provides a resource to create a ci media_concat_template
---

# tencentcloud_ci_media_concat_template

Provides a resource to create a ci media_concat_template

## Example Usage

```hcl
resource "tencentcloud_ci_media_concat_template" "media_concat_template" {
  bucket = "terraform-ci-xxxxxx"
  name   = "concat_templates"
  concat_template {
    concat_fragment {
      url  = "https://terraform-ci-xxxxxx.cos.ap-guangzhou.myqcloud.com/mp4%2Fmp4-test.mp4"
      mode = "Start"
    }
    concat_fragment {
      url  = "https://terraform-ci-xxxxxx.cos.ap-guangzhou.myqcloud.com/mp4%2Fmp4-test.mp4"
      mode = "End"
    }
    audio {
      codec      = "mp3"
      samplerate = ""
      bitrate    = ""
      channels   = ""
    }
    video {
      codec   = "H.264"
      width   = "1280"
      height  = ""
      bitrate = "1000"
      fps     = "25"
      crf     = ""
      remove  = ""
      rotate  = ""
    }
    container {
      format = "mp4"
    }
    audio_mix {
      audio_source = "https://terraform-ci-xxxxxx.cos.ap-guangzhou.myqcloud.com/mp3%2Fnizhan-test.mp3"
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
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) bucket name.
* `concat_template` - (Required, List) stitching template.
* `name` - (Required, String) The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.

The `audio_mix` object supports the following:

* `audio_source` - (Required, String) The media address of the audio track that needs to be mixed.
* `effect_config` - (Optional, List) Mix Fade Configuration.
* `mix_mode` - (Optional, String) Mixing mode Repeat: background sound loop, Once: The background sound is played once.
* `replace` - (Optional, String) Whether to replace the original audio of the Input media file with the mixed audio track media.

The `audio` object supports the following:

* `codec` - (Required, String) Codec format, value aac, mp3.
* `bitrate` - (Optional, String) Original audio bit rate, unit: Kbps, Value range: [8, 1000].
* `channels` - (Optional, String) number of channels- When Codec is set to aac, support 1, 2, 4, 5, 6, 8- When Codec is set to mp3, support 1, 2.
* `samplerate` - (Optional, String) Sampling Rate- Unit: Hz- Optional 11025, 22050, 32000, 44100, 48000, 96000- Different packages, mp3 supports different sampling rates, as shown in the table below.

The `concat_fragment` object supports the following:

* `mode` - (Required, String) node type, `start`, `end`.
* `url` - (Required, String) Splicing object address.

The `concat_template` object supports the following:

* `concat_fragment` - (Required, List) Package format.
* `container` - (Required, List) Only splicing without transcoding.
* `audio_mix` - (Optional, List) mixing parameters.
* `audio` - (Optional, List) audio parameters, the target file does not require Audio information, need to set Audio.Remove to true.
* `video` - (Optional, List) video information, do not upload Video, which is equivalent to deleting video information.

The `container` object supports the following:

* `format` - (Required, String) Container format: mp4, flv, hls, ts, mp3, aac.

The `effect_config` object supports the following:

* `bgm_fade_time` - (Optional, String) bgm transition fade-in duration, support floating point numbers.
* `enable_bgm_fade` - (Optional, String) Enable bgm conversion fade in.
* `enable_end_fadeout` - (Optional, String) enable fade out.
* `enable_start_fadein` - (Optional, String) enable fade in.
* `end_fadeout_time` - (Optional, String) fade out time, greater than 0, support floating point numbers.
* `start_fadein_time` - (Optional, String) Fade in duration, greater than 0, support floating point numbers.

The `video` object supports the following:

* `codec` - (Required, String) Codec format `H.264`.
* `bitrate` - (Optional, String) Original audio bit rate, unit: Kbps, Value range: [8, 1000].
* `crf` - (Optional, String) Bit rate-quality control factor, value range: (0, 51], If Crf is set, the setting of Bitrate will be invalid, When Bitrate is empty, the default is 25.
* `fps` - (Optional, String) Frame rate, value range: (0, 60], Unit: fps.
* `height` - (Optional, String) High, value range: [128, 4096], Unit: px, If only Height is set, Width is calculated according to the original ratio of the video, must be even.
* `remove` - (Optional, String) Whether to delete the source audio stream, the value is true, false.
* `rotate` - (Optional, String) Rotation angle, Value range: [0, 360), Unit: degree.
* `width` - (Optional, String) width, value range: [128, 4096], Unit: px, If only Width is set, Height is calculated according to the original ratio of the video, must be even.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci media_concat_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_concat_template.media_concat_template id=terraform-ci-xxxxxx#t1cb115dfa1fcc414284f83b7c69bcedcf
```

