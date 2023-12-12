Provides a resource to create a mps word_sample

Example Usage

```hcl
resource "tencentcloud_mps_word_sample" "word_sample" {
  usages = ["Recognition.Ocr","Review.Ocr","Review.Asr"]
  keyword = "tf_test_kw_1"
  tags = ["tags_1", "tags_2"]
}
```

Import

mps word_sample can be imported using the id, e.g.

```
terraform import tencentcloud_mps_word_sample.word_sample keyword
```