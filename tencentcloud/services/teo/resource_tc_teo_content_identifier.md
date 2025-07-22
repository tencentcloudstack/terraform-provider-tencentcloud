Provides a resource to create a TEO content identifier

Example Usage

```hcl
resource "tencentcloud_teo_content_identifier" "example" {
  plan_id     = "edgeone-6bzvsgjkfa9g"
  description = "example"
  tags {
    tag_key   = "tagKey"
    tag_value = "tagValue"
  }
}
```

Import

TEO content identifier can be imported using the id, e.g.

```
terraform import tencentcloud_teo_content_identifier.example eocontent-3dy8iyfq8dba
```
