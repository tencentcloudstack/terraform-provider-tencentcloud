Provides a resource to create a css timeshift_template

Example Usage

```hcl
resource "tencentcloud_css_timeshift_template" "timeshift_template" {
    area                   = "Mainland"
    description            = "timeshift template"
    duration               = 604800
    item_duration          = 5
    remove_watermark       = true
    template_name          = "tf-test"
    transcode_template_ids = []
}
```

Import

css timeshift_template can be imported using the id, e.g.

```
terraform import tencentcloud_css_timeshift_template.timeshift_template templateId
```