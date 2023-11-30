Provides a resource to create a ci media_pic_process_template

Example Usage

```hcl
resource "tencentcloud_ci_media_pic_process_template" "media_pic_process_template" {
  bucket = "terraform-ci-xxxxxx"
  name = "pic_process_template"
  pic_process {
		is_pic_info = "true"
		process_rule = "imageMogr2/rotate/90"

  }
}
```

Import

ci media_pic_process_template can be imported using the bucket#templateId, e.g.

```
terraform import tencentcloud_ci_media_pic_process_template.media_pic_process_template terraform-ci-xxxxx#t184a8a26da4674c80bf260c1e34131a65
```