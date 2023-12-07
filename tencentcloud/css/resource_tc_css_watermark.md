Provides a resource to create a css watermark

Example Usage

```hcl
resource "tencentcloud_css_watermark" "watermark" {
  picture_url = "picture_url"
  watermark_name = "watermark_name"
  x_position = 0
  y_position = 0
  width = 0
  height = 0
}

```
Import

css watermark can be imported using the id, e.g.
```
$ terraform import tencentcloud_css_watermark.watermark watermark_id
```