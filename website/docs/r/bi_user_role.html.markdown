---
subcategory: "Business Intelligence(BI)"
layout: "tencentcloud"
page_title: "TencentCloud: tencentcloud_bi_user_role"
sidebar_current: "docs-tencentcloud-resource-bi_user_role"
description: |-
  Provides a resource to create a bi user_role
---

# tencentcloud_bi_user_role

Provides a resource to create a bi user_role

## Example Usage

```hcl
resource "tencentcloud_bi_user_role" "user_role" {
  area_code    = "+83"
  email        = "1055000000@qq.com"
  phone_number = "13470010000"
  role_id_list = [
    10629359,
  ]
  user_id   = "100032767426"
  user_name = "keep-iac-test"
}
```

## Argument Reference

The following arguments are supported:

* `area_code` - (Required, String) Mobile area code(Note: This field may return null, indicating that no valid value can be obtained).
* `email` - (Required, String) E-mail(Note: This field may return null, indicating that no valid value can be obtained).
* `phone_number` - (Required, String) Phone number(Note: This field may return null, indicating that no valid value can be obtained).
* `role_id_list` - (Required, Set: [`Int`]) Role id list.
* `user_id` - (Required, String) User id.
* `user_name` - (Required, String) Username.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the resource.



## Import

bi user_role can be imported using the id, e.g.

```
terraform import tencentcloud_bi_user_role.user_role user_id
```

