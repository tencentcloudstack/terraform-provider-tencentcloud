Provides a resource to create a css snapshot_template

Example Usage

```hcl
resource "tencentcloud_css_snapshot_template" "snapshot_template" {
    cos_app_id        = 1308919341
    cos_bucket        = "keep-bucket"
    cos_region        = "ap-guangzhou"
    description       = "snapshot template"
    height            = 0
    porn_flag         = 0
    snapshot_interval = 2
    template_name     = "tf-snapshot-template"
    width             = 0
}
```

Import

css snapshot_template can be imported using the id, e.g.

```
terraform import tencentcloud_css_snapshot_template.snapshot_template templateId
```