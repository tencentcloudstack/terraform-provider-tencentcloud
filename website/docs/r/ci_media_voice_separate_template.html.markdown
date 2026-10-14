---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_media_voice_separate_template"
sidebar_current: "docs-tencentcloud-resource-ci_media_voice_separate_template"
description: |-
  Provides a resource to create a ci media_voice_separate_template
---

# tencentcloud_ci_media_voice_separate_template

Provides a resource to create a ci media_voice_separate_template

## Example Usage

```hcl
resource "tencentcloud_ci_media_voice_separate_template" "media_voice_separate_template" {
  bucket     = "terraform-ci-xxxxx"
  name       = "voice_separate_template"
  audio_mode = "IsAudio"
  audio_config {
    codec      = "aac"
    samplerate = "44100"
    bitrate    = "128"
    channels   = "4"
  }
}
```

## Argument Reference

The following arguments are supported:

* `audio_config` - (Required, List) audio configuration.
* `audio_mode` - (Required, String) Output audio IsAudio: output human voice, IsBackground: output background sound, AudioAndBackground: output vocal and background sound.
* `bucket` - (Required, String) bucket name.
* `name` - (Required, String) The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.

The `audio_config` object supports the following:

* `codec` - (Required, String) Codec format, value aac, mp3, flac, amr.
* `bitrate` - (Optional, String) Original audio bit rate, unit: Kbps, Value range: [8, 1000].
* `channels` - (Optional, String) number of channels- When Codec is set to aac/flac, support 1, 2, 4, 5, 6, 8- When Codec is set to mp3, support 1, 2- When Codec is set to amr, only 1 is supported.
* `samplerate` - (Optional, String) Sampling Rate- 1: Unit: Hz- 2: Optional 8000, 11025, 22050, 32000, 44100, 48000, 96000- 3: When Codec is set to aac/flac, 8000 is not supported- 4: When Codec is set to mp3, 8000 and 96000 are not supported- 5: When Codec is set to amr, only 8000 is supported.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci media_voice_separate_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_voice_separate_template.media_voice_separate_template terraform-ci-xxxxxx#t1c95566664530460d9bc2b6265feb7c32
```

