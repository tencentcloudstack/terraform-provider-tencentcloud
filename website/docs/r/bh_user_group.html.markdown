---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bh_user_group"
sidebar_current: "docs-tencentcloud-resource-bh_user_group"
description: |-
  Provides a resource to create a BH user group
---

# tencentcloud_bh_user_group

Provides a resource to create a BH user group

## Example Usage

```hcl
resource "tencentcloud_bh_user_group" "example" {
  name = "tf-example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) User group name, maximum length 32 characters.
* `department_id` - (Optional, String) Department ID to which the user group belongs, e.g.: 1.2.3.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.
* `user_group_id` - User group ID.


## Import

BH user group can be imported using the id, e.g.

```
terraform import tencentcloud_bh_user_group.example 92
```

