Provides a resource to create a ci media_video_process_template

Example Usage

```hcl

resource "tencentcloud_ci_media_video_process_template" "media_video_process_template" {
  bucket = "terraform-ci-xxxxxx"
  name = "video_process_template"
  color_enhance {
		enable = "true"
		contrast = ""
		correction = ""
		saturation = ""

  }
  ms_sharpen {
		enable = "false"
		sharpen_level = ""

  }
}
```

Import

ci media_video_process_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_video_process_template.media_video_process_template terraform-ci-xxxxxx#t1d5694d87639a4593a9fd7e9025d26f52
```