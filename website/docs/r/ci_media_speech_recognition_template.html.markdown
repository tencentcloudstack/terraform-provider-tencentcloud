---
subcategory: "Cloud Infinite(CI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_ci_media_speech_recognition_template"
sidebar_current: "docs-tencentcloud-resource-ci_media_speech_recognition_template"
description: |-
  Provides a resource to create a ci media_speech_recognition_template
---

# tencentcloud_ci_media_speech_recognition_template

Provides a resource to create a ci media_speech_recognition_template

## Example Usage

```hcl
resource "tencentcloud_ci_media_speech_recognition_template" "media_speech_recognition_template" {
  bucket = "terraform-ci-1308919341"
  name   = "speech_recognition_template"
  speech_recognition {
    engine_model_type   = "16k_zh"
    channel_num         = "1"
    res_text_format     = "1"
    filter_dirty        = "0"
    filter_modal        = "1"
    convert_num_mode    = "0"
    speaker_diarization = "1"
    speaker_number      = "0"
    filter_punc         = "0"
    output_file_type    = "txt"
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required, String) bucket name.
* `name` - (Required, String) The template name only supports `Chinese`, `English`, `numbers`, `_`, `-` and `*`.
* `speech_recognition` - (Required, List) audio configuration.

The `speech_recognition` object supports the following:

* `channel_num` - (Required, String) Number of voice channels: 1 means mono. EngineModelType supports only mono for non-telephone scenarios, and 2 means dual channels (only 8k_zh engine model supports dual channels, which should correspond to both sides of the call).
* `engine_model_type` - (Required, String) Engine model type, divided into phone scene and non-phone scene, phone scene: 8k_zh: phone 8k Chinese Mandarin general (can be used for dual-channel audio), 8k_zh_s: phone 8k Chinese Mandarin speaker separation (only for monophonic audio), 8k_en: Telephone 8k English; non-telephone scene: 16k_zh: 16k Mandarin Chinese, 16k_zh_video: 16k audio and video field, 16k_en: 16k English, 16k_ca: 16k Cantonese, 16k_ja: 16k Japanese, 16k_zh_edu: Chinese education, 16k_en_edu: English education, 16k_zh_medical: medical, 16k_th: Thai, 16k_zh_dialect: multi-dialect, supports 23 dialects.
* `convert_num_mode` - (Optional, String) Whether to perform intelligent conversion of Arabic numerals (currently supports Mandarin Chinese engine): 0 means no conversion, directly output Chinese numbers, 1 means intelligently convert to Arabic numerals according to the scene, 3 means enable math-related digital conversion, the default value is 0.
* `filter_dirty` - (Optional, String) Whether to filter dirty words (currently supports Mandarin Chinese engine): 0 means not to filter dirty words, 1 means to filter dirty words, 2 means to replace dirty words with *, the default value is 0.
* `filter_modal` - (Optional, String) Whether to pass modal particles (currently supports Mandarin Chinese engine): 0 means not to filter modal particles, 1 means partial filtering, 2 means strict filtering, and the default value is 0.
* `filter_punc` - (Optional, String) Whether to filter punctuation (currently supports Mandarin Chinese engine): 0 means no filtering, 1 means filtering end-of-sentence punctuation, 2 means filtering all punctuation, the default value is 0.
* `output_file_type` - (Optional, String) Output file type, optional txt, srt. The default is txt.
* `res_text_format` - (Optional, String) Recognition result return form: 0 means the recognition result text (including segmented time stamps), 1 is the detailed recognition result at the word level granularity, without punctuation, and includes the speech rate value (a list of word time stamps, generally used to generate subtitle scenes), 2 Detailed recognition results at word-level granularity (including punctuation and speech rate values)..
* `speaker_diarization` - (Optional, String) Whether to enable speaker separation: 0 means not enabled, 1 means enabled (only supports 8k_zh, 16k_zh, 16k_zh_video, monophonic audio), the default value is 0, Note: 8K telephony scenarios suggest using dual-channel to distinguish between the two parties, set ChannelNum=2 is enough, no need to enable speaker separation.
* `speaker_number` - (Optional, String) The number of speakers to be separated (need to be used in conjunction with enabling speaker separation), value range: 0-10, 0 means automatic separation (currently only supports <= 6 people), 1-10 represents the number of specified speakers to be separated. The default value is 0.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

ci media_speech_recognition_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_speech_recognition_template.media_speech_recognition_template terraform-ci-xxxxxx#t1d794430f2f1f4350b11e905ce2c6167e
```

