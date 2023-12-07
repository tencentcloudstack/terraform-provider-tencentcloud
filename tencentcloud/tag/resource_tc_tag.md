Provides a resource to create a tag

Example Usage

```hcl

resource "tencentcloud_tag" "tag" {
	tag_key = "test"
	tag_value = "Terraform"
}

```

Import

tag tag can be imported using the id, e.g.

```
terraform import tencentcloud_tag.tag tag_id
```