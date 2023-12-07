Provides a resource to create a ci media_watermark_template

Example Usage

```hcl
resource "tencentcloud_ci_media_watermark_template" "media_watermark_template" {
  bucket = "terraform-ci-1308919341"
  name = "watermark_template"
  watermark {
		type = "Text"
		pos = "TopRight"
		loc_mode = "Absolute"
		dx = "128"
		dy = "128"
		start_time = "0"
		end_time = "100.5"
		# image {
		# 	url = ""
		# 	mode = ""
		# 	width = ""
		# 	height = ""
		# 	transparency = ""
		# 	background = ""
		# }
		text {
      font_size = "30"
			font_type = "simfang.ttf"
			font_color = "0xF0F8F0"
			transparency = "30"
			text = "watermark-content"
		}
  }
}
```

Import

ci media_watermark_template can be imported using the id, e.g.

```
terraform import tencentcloud_ci_media_watermark_template.media_watermark_template media_watermark_template_id
```