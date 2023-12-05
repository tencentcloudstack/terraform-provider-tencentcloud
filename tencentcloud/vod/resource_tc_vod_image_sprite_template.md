Provide a resource to create a VOD image sprite template.

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
```

Import

VOD image sprite template can be imported using the id, e.g.

```
$ terraform import tencentcloud_vod_image_sprite_template.foo 51156
```