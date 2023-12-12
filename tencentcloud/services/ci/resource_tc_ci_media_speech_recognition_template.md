Provides a resource to create a ci media_speech_recognition_template

Example Usage

```hcl
resource "tencentcloud_ci_media_speech_recognition_template" "media_speech_recognition_template" {
  bucket = "terraform-ci-1308919341"
  name = "speech_recognition_template"
  speech_recognition {
		engine_model_type = "16k_zh"
		channel_num = "1"
		res_text_format = "1"
		filter_dirty = "0"
		filter_modal = "1"
		convert_num_mode = "0"
		speaker_diarization = "1"
		speaker_number = "0"
		filter_punc = "0"
		output_file_type = "txt"
  }
}
```

Import

ci media_speech_recognition_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_speech_recognition_template.media_speech_recognition_template terraform-ci-xxxxxx#t1d794430f2f1f4350b11e905ce2c6167e
```