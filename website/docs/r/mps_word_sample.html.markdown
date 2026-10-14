---
subcategory: "Media Processing Service(MPS)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_mps_word_sample"
sidebar_current: "docs-tencentcloud-resource-mps_word_sample"
description: |-
  Provides a resource to create a mps word_sample
---

# tencentcloud_mps_word_sample

Provides a resource to create a mps word_sample

## Example Usage

```hcl
resource "tencentcloud_mps_word_sample" "word_sample" {
  usages  = ["Recognition.Ocr", "Review.Ocr", "Review.Asr"]
  keyword = "tf_test_kw_1"
  tags    = ["tags_1", "tags_2"]
}
```

## Argument Reference

The following arguments are supported:

* `keyword` - (Required, String) Keyword. Length limit: 20 characters.
* `usages` - (Required, Set: [`String`]) Keyword usage. Valid values: 1.`Recognition.Ocr`: OCR-based content recognition. 2.`Recognition.Asr`: ASR-based content recognition. 3.`Review.Ocr`: OCR-based inappropriate information recognition. 4.`Review.Asr`: ASR-based inappropriate information recognition.
* `tags` - (Optional, Set: [`String`]) Keyword tag. Array length limit: 20 tags. Each tag length limit: 128 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateWordSamples](https://cloud.tencent.com/document/api/862/39444) - This API is used to create keyword samples in batches for video processing operations such as content recognition and inappropriate information detection with the help of the OCR and ASR technologies.
* Read: [DescribeWordSamples](https://cloud.tencent.com/document/api/862/39440) - This API is used to perform paged queries of keyword sample information by use case, keyword, and tag.
* Update: [ModifyWordSample](https://cloud.tencent.com/document/api/862/39438) - This API is used to modify the use case and tag of a keyword. The keyword itself cannot be modified, but you can delete it and create another one if needed.
* Delete: [DeleteWordSamples](https://cloud.tencent.com/document/api/862/39442) - This API is used to delete keyword samples in batches.

## Import

mps word_sample can be imported using the id, e.g.

```
terraform import tencentcloud_mps_word_sample.word_sample keyword
```

