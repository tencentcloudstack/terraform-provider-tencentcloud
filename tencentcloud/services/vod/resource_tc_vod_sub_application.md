Provide a resource to create a VOD sub application.

Example Usage

### Basic Usage

```hcl
resource "tencentcloud_vod_sub_application" "example" {
  name        = "tf-example"
  status      = "On"
  description = "this is sub application"
}
```

### Tags Update Example

```hcl
resource "tencentcloud_vod_sub_application" "example" {
  name           = "tf-example"
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

VOD sub application can be imported using the name and id separated by `name#sub_app_id`, e.g.

```
terraform import tencentcloud_vod_sub_application.example tf-example#1500066377
```
