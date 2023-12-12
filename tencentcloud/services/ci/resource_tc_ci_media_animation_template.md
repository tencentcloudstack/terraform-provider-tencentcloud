Provides a resource to create a ci media_animation_template

Example Usage

```hcl
resource "tencentcloud_ci_media_animation_template" "media_animation_template" {
  bucket = "terraform-ci-1308919341"
  name = "animation_template-002"
  container {
		format = "gif"
  }
  video {
		codec = "gif"
		width = "1280"
		height = ""
		fps = "20"
		animate_only_keep_key_frame = "true"
		animate_time_interval_of_frame = ""
		animate_frames_per_second = ""
		quality = ""

  }
  time_interval {
		start = "0"
		duration = "60"

  }
}
```