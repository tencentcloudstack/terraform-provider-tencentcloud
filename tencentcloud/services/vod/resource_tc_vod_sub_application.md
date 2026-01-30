Provide a resource to create a VOD sub application.

Example Usage

### Basic Usage

```hcl
resource "tencentcloud_vod_sub_application" "foo" {
  name        = "foo"
  status      = "On"
  description = "this is sub application"
}
```

### Tags Update Example

```hcl
resource "tencentcloud_vod_sub_application" "with_tags" {
  name           = "my-app-with-tags"
  status         = "On"
  description    = "Sub application with updatable tags"
  
  tags = {
    "team"        = "media"
    "environment" = "production"
  }
}

# Tags can be updated without recreating the resource
# Modify the tags map to add, update, or remove tags
```

Import

VOD sub application can be imported using the name and id separated by `#`, e.g.

```
$ terraform import tencentcloud_vod_sub_application.foo name#id
```
