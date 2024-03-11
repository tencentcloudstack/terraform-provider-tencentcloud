Provides a resource to create a vod watermark template

Example Usage

```hcl
resource  "tencentcloud_vod_sub_application" "sub_application" {
	name = "watermarkTemplateSubApplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_watermark_template" "watermark_template" {
	type = "image"
	sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
	name = "myImageWatermark"
	comment = "a png watermark"
	coordinate_origin = "TopLeft"
	x_pos = "10%"
	y_pos = "10%"
	image_template {
		image_content = filebase64("xxx.png")
		width = "10%"
		height = "10px"
	}
}
```

Import

vod watermark template can be imported using the id, e.g.

```
terraform import tencentcloud_vod_watermark_template.watermark_template $subAppId#$templateId
```