Use this data source to query detailed information of VOD image sprite templates.

Example Usage

```hcl
resource "tencentcloud_vod_image_sprite_template" "foo" {
  sample_type         = "Percent"
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

data "tencentcloud_vod_image_sprite_templates" "foo" {
  type       = "Custom"
  definition = tencentcloud_vod_image_sprite_template.foo.id
}
```