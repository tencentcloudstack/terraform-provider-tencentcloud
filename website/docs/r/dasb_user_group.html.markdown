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
  name = "tf_example"
}
```

### Or

```hcl
resource "tencentcloud_dasb_user_group" "example" {
  name          = "tf_example"
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


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateUserGroup](https://www.tencentcloud.com/search/?wd=CreateUserGroup&type=document&lang=en)
* Read: [DescribeUserGroups](https://www.tencentcloud.com/search/?wd=DescribeUserGroups&type=document&lang=en)
* Update: [ModifyUserGroup](https://www.tencentcloud.com/search/?wd=ModifyUserGroup&type=document&lang=en)
* Delete: [DeleteUserGroups](https://www.tencentcloud.com/search/?wd=DeleteUserGroups&type=document&lang=en)

## Import

dasb user_group can be imported using the id, e.g.

```
terraform import tencentcloud_dasb_user_group.example 16
```

