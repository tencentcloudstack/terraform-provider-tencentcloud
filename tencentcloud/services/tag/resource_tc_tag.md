Provides a resource to create a tag

Example Usage

```hcl
resource "tencentcloud_tag" "example" {
  tag_key   = "tagKey"
  tag_value = "tagValue"
}
```

Import

tag can be imported using the id, e.g.

```
terraform import tencentcloud_tag.example tagKey#tagValue
```