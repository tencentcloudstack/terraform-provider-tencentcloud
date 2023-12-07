Provides a resource to create a mps ai_recognition_template

Example Usage

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
    default_library_label_set     = [
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

Import

mps ai_recognition_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_ai_recognition_template.ai_recognition_template ai_recognition_template_id
```