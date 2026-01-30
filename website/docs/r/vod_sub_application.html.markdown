---
subcategory: "Video on Demand(VOD)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_vod_sub_application"
sidebar_current: "docs-tencentcloud-resource-vod_sub_application"
description: |-
  Provide a resource to create a VOD sub application.
---

# tencentcloud_vod_sub_application

Provide a resource to create a VOD sub application.

## Example Usage

### ### Basic Usage

```hcl
resource "tencentcloud_vod_sub_application" "foo" {
  name        = "foo"
  status      = "On"
  description = "this is sub application"
}
```

### ### Tags Update Example

```hcl
resource "tencentcloud_vod_sub_application" "with_tags" {
  name        = "my-app-with-tags"
  status      = "On"
  description = "Sub application with updatable tags"

  tags = {
    "team"        = "media"
    "environment" = "production"
  }
}

# Tags can be updated without recreating the resource
# Modify the tags map to add, update, or remove tags
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Sub application name, which can contain up to 64 letters, digits, underscores, and hyphens (such as test_ABC-123) and must be unique under a user.
* `status` - (Required, String) Sub appliaction status.
* `description` - (Optional, String) Sub application description.
* `tags` - (Optional, Map) Tag key-value pairs for resource management. Maximum 10 tags.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time when the sub application was created.


## Import

VOD sub application can be imported using the name and id separated by `#`, e.g.

```
$ terraform import tencentcloud_vod_sub_application.foo name#id
```

