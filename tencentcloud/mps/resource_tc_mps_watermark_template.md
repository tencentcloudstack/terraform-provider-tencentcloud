Provides a resource to create a mps watermark_template

Example Usage

```hcl
resource "tencentcloud_mps_watermark_template" "watermark_template" {
  coordinate_origin = "TopLeft"
  name              = "xZxasd"
  type              = "image"
  x_pos             = "12%"
  y_pos             = "21%"

  image_template {
    height        = "17px"
    image_content = filebase64("./logo.png")
    repeat_type   = "repeat"
    width         = "12px"
  }
}
```

Import

mps watermark_template can be imported using the id, e.g.

```
terraform import tencentcloud_mps_watermark_template.watermark_template watermark_template_id
```