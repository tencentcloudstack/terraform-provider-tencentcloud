---
subcategory: "Bastion Host(BH)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_dasb_user_group"
sidebar_current: "docs-tencentcloud-resource-dasb_user_group"
description: |-
  Provides a resource to create a dasb user_group
---

# tencentcloud_dasb_user_group

Provides a resource to create a dasb user_group

## Example Usage

```hcl
resource "tencentcloud_dasb_user_group" "example" {
  name = "tf_example_update"
}
```

### Or

```hcl
resource "tencentcloud_dasb_user_group" "example" {
  name          = "tf_example_update"
  department_id = "1.2"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) User group name, maximum length 32 characters.
* `department_id` - (Optional, String) ID of the department to which the user group belongs, such as: 1.2.3.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

dasb user_group can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_user_group.example 16
```

