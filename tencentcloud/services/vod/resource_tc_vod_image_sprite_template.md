Provide a resource to create a VOD image sprite template.

Example Usage

```hcl
resource  "tencentcloud_vod_sub_application" "sub_application" {
	name = "image-sprite-subapplication"
	status = "On"
	description = "this is sub application"
}

resource "tencentcloud_vod_image_sprite_template" "foo" {
  sample_type         = "Percent"
  sub_app_id = tonumber(split("#", tencentcloud_vod_sub_application.sub_application.id)[1])
  sample_interval     = 10
  row_count           = 3
  column_count        = 3
  name                = "tf-sprite"
  comment             = "test"
  fill_type           = "stretch"
  width               = 128
  height              = 128
  resolution_adaptive = false
}
```

Import

VOD image sprite template can be imported using the id($subAppId#$templateId), e.g.

```
$ terraform import tencentcloud_vod_image_sprite_template.foo $subAppId#$templateId
```