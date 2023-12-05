Provides a resource to create a ci media_super_resolution_template

Example Usage

```hcl
resource "tencentcloud_ci_media_super_resolution_template" "media_super_resolution_template" {
  bucket = "terraform-ci-1308919341"
  name = "super_resolution_template"
  resolution = "sdtohd"
  enable_scale_up = "true"
  version = "Enhance"
}
```

Import

ci media_super_resolution_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_super_resolution_template.media_super_resolution_template terraform-ci-xxxxxx#t1d707eb2be3294e22b47123894f85cb8f
```