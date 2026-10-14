---
subcategory: "Business Intelligence(BI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bi_project_user_role"
sidebar_current: "docs-tencentcloud-resource-bi_project_user_role"
description: |-
  Provides a resource to create a bi project_user_role
---

# tencentcloud_bi_project_user_role

Provides a resource to create a bi project_user_role

~> **NOTE:** You cannot use `tencentcloud_bi_user_role` and `tencentcloud_bi_project_user_role` at the same time to modify the `phone_number` and `email` of the same user.

## Example Usage

```hcl
resource "tencentcloud_bi_project_user_role" "project_user_role" {
  area_code    = "+86"
  project_id   = 11015030
  role_id_list = [10629453]
  email        = "123456@qq.com"
  phone_number = "13130001000"
  user_id      = "100024664626"
  user_name    = "keep-cam-user"
}
```

## Argument Reference

The following arguments are supported:

* `area_code` - (Required, String) Mobile area code(Note: This field may return null, indicating that no valid value can be obtained).
* `email` - (Required, String) E-mail(Note: This field may return null, indicating that no valid value can be obtained).
* `phone_number` - (Required, String) Phone number(Note: This field may return null, indicating that no valid value can be obtained).
* `user_id` - (Required, String) User id.
* `user_name` - (Required, String) Username.
* `project_id` - (Optional, Int) Project id.
* `role_id_list` - (Optional, Set: [`Int`]) Role id list.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.


## API Reference

This resource interacts with the following TencentCloud APIs:

* Create: [CreateUserRoleProject](https://cloud.tencent.com/document/api/590/98847) - This API is used to create a user role in the project.
* Read: [DescribeUserRoleProjectList](https://cloud.tencent.com/document/api/590/100532) - This API is used to obtain the user role list in the project.
* Update: [ModifyUserRoleProject](https://cloud.tencent.com/document/api/590/98842) - This API is used to modify the user role info in the project.
* Delete: [DeleteUserRoleProject](https://cloud.tencent.com/document/api/590/98845) - This API is used to delete a user role in the project.

## Import

bi project_user_role can be imported using the id, e.g.

```
terraform import tencentcloud_bi_project_user_role.project_user_role projectId#userId
```

