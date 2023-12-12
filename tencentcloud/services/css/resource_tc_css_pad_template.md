Provides a resource to create a css pad_template

Example Usage

```hcl
resource "tencentcloud_css_pad_template" "pad_template" {
  description   = "pad template"
  max_duration  = 120000
  template_name = "tf-pad"
  type          = 1
  url           = "https://livewatermark-1251132611.cos.ap-guangzhou.myqcloud.com/1308919341/watermark_img_1698736540399_1441698123618_.pic.jpg"
  wait_duration = 2000
}
```

Import

css pad_template can be imported using the id, e.g.

```
terraform import tencentcloud_css_pad_template.pad_template templateId
```