Provides a resource to create a mps image_sprite_template

Example Usage

```hcl
resource "tencentcloud_mps_image_sprite_template" "image_sprite_template" {
  column_count        = 10
  fill_type           = "stretch"
  format              = "jpg"
  height              = 143
  name                = "terraform-test"
  resolution_adaptive = "open"
  row_count           = 10
  sample_interval     = 10
  sample_type         = "Time"
  width               = 182
}
```

Import

mps image_sprite_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_image_sprite_template.image_sprite_template image_sprite_template_id
```