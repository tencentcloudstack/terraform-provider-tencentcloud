---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_ai_recognition_template"
sidebar_current: "docs-tencentcloud-resource-mps_ai_recognition_template"
description: |-
  Provides a resource to create a mps ai_recognition_template
---

# tencentcloud_mps_ai_recognition_template

Provides a resource to create a mps ai_recognition_template

## Example Usage

```hcl
resource "tencentcloud_mps_ai_recognition_template" "ai_recognition_template" {
  name = "terraform-test"

  asr_full_text_configure {
    switch = "OFF"
  }

  asr_words_configure {
    label_set = []
    switch    = "OFF"
  }

  face_configure {
    default_library_label_set = [
      "entertainment",
      "sport",
    ]
    face_library                  = "All"
    score                         = 85
    switch                        = "ON"
    user_define_library_label_set = []
  }

  ocr_full_text_configure {
    switch = "OFF"
  }

  ocr_words_configure {
    label_set = []
    switch    = "OFF"
  }
}
```

## Argument Reference

The following arguments are supported:

* `asr_full_text_configure` - (Optional, List) Asr full text recognition control parameters.
* `asr_words_configure` - (Optional, List) Asr word recognition control parameters.
* `comment` - (Optional, String) Ai recognition template description information, length limit: 256 characters.
* `face_configure` - (Optional, List) Face recognition control parameters.
* `name` - (Optional, String) Ai recognition template name, length limit: 64 characters.
* `ocr_full_text_configure` - (Optional, List) Ocr full text control parameters.
* `ocr_words_configure` - (Optional, List) Ocr words recognition control parameters.

The `asr_full_text_configure` object supports the following:

* `switch` - (Required, String) Asr full text recognition task switch, optional value:ON/OFF.
* `subtitle_format` - (Optional, String) Generated subtitle file format, if left blank or blank string means no subtitle file will be generated, optional value:vtt: Generate WebVTT subtitle files.

The `asr_words_configure` object supports the following:

* `switch` - (Required, String) Asr word recognition task switch, optional value:ON/OFF.
* `label_set` - (Optional, Set) Keyword filter label, specify the label of the keyword to be returned. If not filled or empty, all results will be returned.The maximum number of tags is 10, and the length of each tag is up to 16 characters.

The `face_configure` object supports the following:

* `switch` - (Required, String) Ai face recognition task switch, optional value:ON/OFF.
* `default_library_label_set` - (Optional, Set) Default face filter tag, specify the tag of the default face that needs to be returned. If not filled or empty, all default face results will be returned. Label optional value:entertainment, sport, politician.
* `face_library` - (Optional, String) Face library selection, optional value:Default, UserDefine, AllDefault value: All, use the system default face library and user-defined face library.
* `score` - (Optional, Float64) Face recognition filter score, when the recognition result reaches the score above, the recognition result will be returned. The default is 95 points. Value range: 0 - 100.
* `user_define_library_label_set` - (Optional, Set) User-defined face filter tag, specify the tag of the user-defined face that needs to be returned. If not filled or empty, all custom face results will be returned.The maximum number of tags is 100, and the length of each tag is up to 16 characters.

The `ocr_full_text_configure` object supports the following:

* `switch` - (Required, String) Ocr full text recognition task switch, optional value:ON/OFF.

The `ocr_words_configure` object supports the following:

* `switch` - (Required, String) Ocr words recognition task switch, optional value:ON/OFF.
* `label_set` - (Optional, Set) Keyword filter label, specify the label of the keyword to be returned. If not filled or empty, all results will be returned.The maximum number of tags is 10, and the length of each tag is up to 16 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps ai_recognition_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_ai_recognition_template.ai_recognition_template ai_recognition_template_id
```

