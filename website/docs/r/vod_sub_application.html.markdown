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

```hcl
resource "tencentcloud_vod_sub_application" "foo" {
  name        = "foo"
  status      = "On"
  description = "this is sub application"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Sub application name, which can contain up to 64 letters, digits, underscores, and hyphens (such as test_ABC-123) and must be unique under a user.
* `status` - (Required, String) Sub appliaction status.
* `description` - (Optional, String) Sub application description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - The time when the sub application was created.


## Import

VOD super player config can be imported using the name+, e.g.

```
$ terraform import tencentcloud_vod_sub_application.foo name+"#"+id
```

