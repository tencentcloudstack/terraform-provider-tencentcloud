Provides a resource to create a mps animated_graphics_template

Example Usage

```hcl
resource "tencentcloud_mps_animated_graphics_template" "animated_graphics_template" {
  format              = "gif"
  fps                 = 20
  height              = 130
  name                = "terraform-test"
  quality             = 75
  resolution_adaptive = "open"
  width               = 140
}
```

Import

mps animated_graphics_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_animated_graphics_template.animated_graphics_template animated_graphics_template_id
```