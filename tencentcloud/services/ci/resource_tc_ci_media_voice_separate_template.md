Provides a resource to create a ci media_voice_separate_template

Example Usage

```hcl
resource "tencentcloud_ci_media_voice_separate_template" "media_voice_separate_template" {
  bucket = "terraform-ci-xxxxx"
  name = "voice_separate_template"
  audio_mode = "IsAudio"
  audio_config {
		codec = "aac"
		samplerate = "44100"
		bitrate = "128"
		channels = "4"
  }
}
```

Import

ci media_voice_separate_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_voice_separate_template.media_voice_separate_template terraform-ci-xxxxxx#t1c95566664530460d9bc2b6265feb7c32
```