---
subcategory: "Tencent Cloud Organization (TCO)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_identity_center_group"
sidebar_current: "docs-tencentcloud-resource-identity_center_group"
description: |-
  Provides a resource to create an identity center group
---

# tencentcloud_identity_center_group

Provides a resource to create an identity center group

## Example Usage

```hcl
resource "tencentcloud_identity_center_group" "identity_center_group" {
  zone_id     = "z-xxxxxx"
  group_name  = "test-group"
  description = "test"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - (Required, String) The name of the user group. Format: Allow English letters, numbers and special characters-. Length: Maximum 128 characters.
* `zone_id` - (Required, String) Zone id.
* `description` - (Optional, String) A description of the user group. Length: Maximum 1024 characters.
* `group_type` - (Optional, String) Type of user group. `Manual`: manual creation, `Synchronized`: external import.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `create_time` - Creation time of the user group.
* `group_id` - ID of the user group.
* `member_count` - Number of team members.
* `update_time` - Modification time for the user group.

## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateGroup](https://www.tencentcloud.com/search/?wd=CreateGroup&type=document&lang=en)
* Read: [GetGroup](https://www.tencentcloud.com/search/?wd=GetGroup&type=document&lang=en)
* Update: [UpdateGroup](https://www.tencentcloud.com/search/?wd=UpdateGroup&type=document&lang=en)
* Delete: [DeleteGroup](https://www.tencentcloud.com/search/?wd=DeleteGroup&type=document&lang=en)

## Import

tencentcloud_identity_center_group can be imported using the id, e.g.

```
terraform import tencentcloud_identity_center_group.identity_center_group ${zoneId}#${groupId}
```

