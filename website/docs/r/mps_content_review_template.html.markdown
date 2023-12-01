---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_content_review_template"
sidebar_current: "docs-tencentcloud-resource-mps_content_review_template"
description: |-
  Provides a resource to create a mps content_review_template
---

# tencentcloud_mps_content_review_template

Provides a resource to create a mps content_review_template

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `comment` - (Optional, String) Content review template description information, length limit: 256 characters.
* `name` - (Optional, String) Content review template name, length limit: 64 characters.
* `political_configure` - (Optional, List) Political control parameters.
* `porn_configure` - (Optional, List) Control parameters for porn image.
* `prohibited_configure` - (Optional, List) Prohibited control parameters. Prohibited content includes:abuse, drug-related violations.Note: this parameter is not yet supported.
* `terrorism_configure` - (Optional, List) Control parameters for unsafe information.
* `user_define_configure` - (Optional, List) User-Defined Content Moderation Control Parameters.

The `asr_review_info` object supports the following:

* `switch` - (Required, String) Political asr task switch, optional value:ON/OFF.
* `block_confidence` - (Optional, Int) The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.
* `review_confidence` - (Optional, Int) The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.

The `asr_review_info` object supports the following:

* `switch` - (Required, String) User-defined asr review task switch, optional value:ON/OFF.
* `block_confidence` - (Optional, Int) The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.
* `label_set` - (Optional, Set) User-defined asr tags, the review result contains the selected tag and returns the result, if the filter tag is empty, all review results are returned. If you want to use the tag filtering function, when adding a asr library, you need to add the corresponding character tag.The maximum number of tags is 10, and the length of each tag is up to 16 characters.
* `review_confidence` - (Optional, Int) The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.

The `asr_review_info` object supports the following:

* `switch` - (Required, String) Voice Prohibition task switch, optional value:ON/OFF.
* `block_confidence` - (Optional, Int) The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.
* `review_confidence` - (Optional, Int) The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.

The `asr_review_info` object supports the following:

* `switch` - (Required, String) Voice pornography task switch, optional value:ON/OFF.
* `block_confidence` - (Optional, Int) The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.
* `review_confidence` - (Optional, Int) The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.

The `face_review_info` object supports the following:

* `switch` - (Required, String) User-defined face review task switch, optional value:ON/OFF.
* `block_confidence` - (Optional, Int) The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.
* `label_set` - (Optional, Set) User-defined face review tags, the review result contains the selected tag and returns the result, if the filter tag is empty, all review results are returned. If you want to use the tag filtering function, when adding a face library, you need to add the corresponding character tag.The maximum number of tags is 10, and the length of each tag is up to 16 characters.
* `review_confidence` - (Optional, Int) The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.

The `img_review_info` object supports the following:

* `switch` - (Required, String) Political image task switch, optional value:ON/OFF.
* `block_confidence` - (Optional, Int) The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 97 points. Value range: 0~100.
* `label_set` - (Optional, Set) Political image filter tag, if the review result contains the selected tag, the result will be returned, if the filter tag is empty, all the review results will be returned, the optional value is:violation_photo, politician, entertainment, sport, entrepreneur, scholar, celebrity, military.
* `review_confidence` - (Optional, Int) The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 95 points. Value range: 0~100.

The `img_review_info` object supports the following:

* `switch` - (Required, String) Porn screen task switch, optional value:ON/OFF.
* `block_confidence` - (Optional, Int) The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 90 points. Value range: 0~100.
* `label_set` - (Optional, Set) Porn image filter label, if the review result contains the selected label, the result will be returned. If the filter label is empty, all the review results will be returned. The optional value is:porn, vulgar, intimacy, sexy.
* `review_confidence` - (Optional, Int) The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 0. Value range: 0~100.

The `img_review_info` object supports the following:

* `switch` - (Required, String) Terrorism image task switch, optional value:ON/OFF.
* `block_confidence` - (Optional, Int) The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 90 points. Value range: 0~100.
* `label_set` - (Optional, Set) Terrorism image filter tag, if the review result contains the selected tag, the result will be returned, if the filter tag is empty, all the review results will be returned, the optional value is:guns, crowd, bloody, police, banners, militant, explosion, terrorists, scenario.
* `review_confidence` - (Optional, Int) The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 80 points. Value range: 0~100.

The `ocr_review_info` object supports the following:

* `switch` - (Required, String) Ocr Prohibition task switch, optional value:ON/OFF.
* `block_confidence` - (Optional, Int) The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.
* `review_confidence` - (Optional, Int) The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.

The `ocr_review_info` object supports the following:

* `switch` - (Required, String) Ocr pornography task switch, optional value:ON/OFF.
* `block_confidence` - (Optional, Int) The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.
* `review_confidence` - (Optional, Int) The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.

The `ocr_review_info` object supports the following:

* `switch` - (Required, String) Ocr terrorism image task switch, optional value:ON/OFF.
* `block_confidence` - (Optional, Int) The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.
* `review_confidence` - (Optional, Int) The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.

The `ocr_review_info` object supports the following:

* `switch` - (Required, String) Political ocr task switch, optional value:ON/OFF.
* `block_confidence` - (Optional, Int) The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.
* `review_confidence` - (Optional, Int) The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.

The `ocr_review_info` object supports the following:

* `switch` - (Required, String) User-defined ocr text review task switch, optional value:ON/OFF.
* `block_confidence` - (Optional, Int) The score threshold for judging suspected violations. When the smart review reaches the score above, it is considered suspected violations. If it is not filled, the default is 100 points. Value range: 0~100.
* `label_set` - (Optional, Set) User-defined ocr tags, the review result contains the selected tag and returns the result, if the filter tag is empty, all review results are returned. If you want to use the tag filtering function, when adding a ocr library, you need to add the corresponding character tag.The maximum number of tags is 10, and the length of each tag is up to 16 characters.
* `review_confidence` - (Optional, Int) The score threshold for judging whether manual review is required for violations. When the intelligent review reaches the score above, it is considered that manual review is required. If it is not filled, the default is 75 points. Value range: 0~100.

The `political_configure` object supports the following:

* `asr_review_info` - (Optional, List) Political asr control parameters.
* `img_review_info` - (Optional, List) Political image control parameters.
* `ocr_review_info` - (Optional, List) Political ocr control parameters.

The `porn_configure` object supports the following:

* `asr_review_info` - (Optional, List) Voice pornography control parameters.
* `img_review_info` - (Optional, List) Porn image Identification Control Parameters.
* `ocr_review_info` - (Optional, List) Ocr pornography control parameters.

The `prohibited_configure` object supports the following:

* `asr_review_info` - (Optional, List) Voice Prohibition Control Parameters.Note: This field may return null, indicating that no valid value can be obtained.
* `ocr_review_info` - (Optional, List) Ocr Prohibition Control Parameters.Note: This field may return null, indicating that no valid value can be obtained.

The `terrorism_configure` object supports the following:

* `ocr_review_info` - (Required, List) Ocr terrorism task Control Parameters.
* `img_review_info` - (Optional, List) Terrorism image task control parameters.

The `user_define_configure` object supports the following:

* `asr_review_info` - (Optional, List) User-defined asr text review control parameters.
* `face_review_info` - (Optional, List) User-defined face review control parameters.
* `ocr_review_info` - (Optional, List) User-defined ocr text review control parameters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

mps content_review_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_content_review_template.content_review_template definition
```

