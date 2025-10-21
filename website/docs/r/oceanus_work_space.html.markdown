---
subcategory: "Oceanus"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_oceanus_work_space"
sidebar_current: "docs-tencentcloud-resource-oceanus_work_space"
description: |-
  Provides a resource to create a oceanus work_space
---

# tencentcloud_oceanus_work_space

Provides a resource to create a oceanus work_space

## Example Usage

```hcl
resource "tencentcloud_oceanus_work_space" "example" {
  work_space_name = "tf_example"
  description     = "example description."
}
```

## Argument Reference

The following arguments are supported:

* `work_space_name` - (Required, String) Workspace name.
* `description` - (Optional, String) Workspace description.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `app_id` - User APPID.
* `create_time` - Create time.
* `creator_uin` - Creator UIN.
* `jobs_count` - Number of Jobs.
* `owner_uin` - Owner UIN.
* `role_auth_count` - Number of workspace members.
* `serial_id` - Serial ID.
* `status` - Workspace status.
* `update_time` - Update time.
* `work_space_id` - Workspace ID.


## Import

oceanus work_space can be imported using the id, e.g.

```
terraform import tencentcloud_oceanus_work_space.example space-0dan3yux#tf_example
```

