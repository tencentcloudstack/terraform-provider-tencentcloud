Provides a resource to create a ci media_smart_cover_template

Example Usage

```hcl
resource "tencentcloud_ci_media_smart_cover_template" "media_smart_cover_template" {
  bucket = "terraform-ci-xxxxxx"
  name = "smart_cover_template"
  smart_cover {
		format = "jpg"
		width = "1280"
		height = "960"
		count = "10"
		delete_duplicates = "true"
  }
}
```

Import

ci media_smart_cover_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_smart_cover_template.media_smart_cover_template terraform-ci-xxxxxx#t1ede83acc305e423799d638044d859fb7
```