Provides a resource to create a ci media_tts_template

Example Usage

```hcl
resource "tencentcloud_ci_media_tts_template" "media_tts_template" {
  bucket = "terraform-ci-xxxxxx"
  name = "tts_template"
  mode = "Asyc"
  codec = "pcm"
  voice_type = "ruxue"
  volume = "0"
  speed = "100"
}
```

Import

ci media_tts_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_tts_template.media_tts_template terraform-ci-xxxxxx#t1ed421df8bd2140b6b73474f70f99b0f8
```