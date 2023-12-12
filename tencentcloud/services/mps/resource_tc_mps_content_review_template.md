Provides a resource to create a mps content_review_template

Example Usage

```hcl
resource "tencentcloud_mps_content_review_template" "template" {
  name    = "tf_test_content_review_temp"
  comment = "tf test content review temp"
  porn_configure {
    img_review_info {
      switch            = "ON"
      label_set         = ["porn", "vulgar"]
      block_confidence  = 60
      review_confidence = 100
    }
    asr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }

  }
  terrorism_configure {
    img_review_info {
      switch            = "ON"
      label_set         = ["guns", "crowd"]
      block_confidence  = 60
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }

  }
  political_configure {
    img_review_info {
      switch            = "ON"
      label_set         = ["violation_photo", "politician"]
      block_confidence  = 60
      review_confidence = 100
    }
    asr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }

  }
  prohibited_configure {
    asr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "ON"
      block_confidence  = 60
      review_confidence = 100
    }

  }
  user_define_configure {
    face_review_info {
      switch            = "ON"
      label_set         = ["FACE_1", "FACE_2"]
      block_confidence  = 60
      review_confidence = 100
    }
    asr_review_info {
      switch            = "ON"
      label_set         = ["VOICE_1", "VOICE_2"]
      block_confidence  = 60
      review_confidence = 100
    }
    ocr_review_info {
      switch            = "ON"
      label_set         = ["VIDEO_1", "VIDEO_2"]
      block_confidence  = 60
      review_confidence = 100
    }
  }
}
```

Import

mps content_review_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_content_review_template.content_review_template definition
```