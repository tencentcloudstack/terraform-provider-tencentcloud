Manage original image protection functionality

Example Usage

```hcl

resource "tencentcloud_ci_original_image_protection" "foo" {
	bucket = "examplebucket-1250000000"
	status = "on"
}

```

Import

Resource original image protection can be imported using the id, e.g.

```
$ terraform import tencentcloud_ci_original_image_protection.example examplebucket-1250000000
```